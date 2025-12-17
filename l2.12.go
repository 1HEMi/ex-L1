package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type intFlag struct {
	val int
	set bool
}

func (f *intFlag) Set(s string) error {
	v, err := strconvAtoi(s)
	if err != nil {
		return err
	}
	f.val = v
	f.set = true
	return nil
}
func (f *intFlag) String() string { return fmt.Sprintf("%d", f.val) }
func (f *intFlag) Get() any       { return f.val }

func strconvAtoi(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty number")
	}
	sign := 1
	i := 0
	if s[0] == '+' {
		i++
	} else if s[0] == '-' {
		sign = -1
		i++
	}
	if i >= len(s) {
		return 0, errors.New("invalid number")
	}
	n := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, errors.New("invalid number")
		}
		n = n*10 + int(c-'0')
	}
	return sign * n, nil
}

type options struct {
	A, B, C intFlag
	count   bool
	icase   bool
	invert  bool
	fixed   bool
	lineno  bool
}

type line struct {
	no   int
	text string
}

func main() {
	var o options

	flag.Var(&o.A, "A", "print N lines of trailing context after matching lines")
	flag.Var(&o.B, "B", "print N lines of leading context before matching lines")
	flag.Var(&o.C, "C", "print N lines of output context (equiv. -A N -B N)")
	flag.BoolVar(&o.count, "c", false, "print only a count of matching lines")
	flag.BoolVar(&o.icase, "i", false, "ignore case distinctions")
	flag.BoolVar(&o.invert, "v", false, "invert the sense of matching")
	flag.BoolVar(&o.fixed, "F", false, "PATTERN is a fixed string (not a regex)")
	flag.BoolVar(&o.lineno, "n", false, "print line number with output lines")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <pattern> [file|-]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		dieUsage("missing pattern")
	}
	pattern := args[0]

	if o.C.set {
		if !o.A.set {
			o.A.val = o.C.val
		}
		if !o.B.set {
			o.B.val = o.C.val
		}
	}

	if o.A.val < 0 || o.B.val < 0 || o.C.val < 0 {
		dieUsage("context values must be >= 0")
	}

	var in io.Reader = os.Stdin
	if len(args) >= 2 && args[1] != "-" {
		f, err := os.Open(args[1])
		if err != nil {
			dieErr(err, 1)
		}
		defer f.Close()
		in = f
	}

	matchFn, err := buildMatcher(pattern, o.fixed, o.icase)
	if err != nil {
		dieErr(err, 2)
	}

	if o.count {
		n, err := countMatches(in, matchFn, o.invert)
		if err != nil {
			dieErr(err, 1)
		}
		fmt.Println(n)
		return
	}

	if err := streamGrep(in, matchFn, o); err != nil {
		dieErr(err, 1)
	}
}

func dieUsage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	flag.Usage()
	os.Exit(2)
}

func dieErr(err error, code int) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(code)
}

func buildMatcher(pattern string, fixed bool, icase bool) (func(string) bool, error) {
	if fixed {
		if icase {
			p := strings.ToLower(pattern)
			return func(s string) bool {
				return strings.Contains(strings.ToLower(s), p)
			}, nil
		}
		return func(s string) bool {
			return strings.Contains(s, pattern)
		}, nil
	}

	if icase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return re.MatchString, nil
}

func countMatches(r io.Reader, matchFn func(string) bool, invert bool) (int, error) {
	sc := newScanner(r)
	cnt := 0
	for sc.Scan() {
		m := matchFn(sc.Text())
		if invert {
			m = !m
		}
		if m {
			cnt++
		}
	}
	return cnt, sc.Err()
}

func streamGrep(r io.Reader, matchFn func(string) bool, o options) error {
	beforeN := o.B.val
	afterN := o.A.val

	sc := newScanner(r)

	var beforeBuf []line
	if beforeN > 0 {
		beforeBuf = make([]line, 0, beforeN)
	}

	afterRemain := 0

	printedAny := false
	lastPrintedNo := 0

	emitSepIfGap := func(nextNo int) {
		if (beforeN == 0 && afterN == 0) || !printedAny {
			return
		}
		if nextNo > lastPrintedNo+1 {
			fmt.Println("--")
		}
	}

	printLine := func(ln line, isMatch bool) {
		emitSepIfGap(ln.no)
		if o.lineno {
			if isMatch {
				fmt.Printf("%d:%s\n", ln.no, ln.text)
			} else {
				fmt.Printf("%d-%s\n", ln.no, ln.text)
			}
		} else {
			fmt.Println(ln.text)
		}
		printedAny = true
		lastPrintedNo = ln.no
	}

	pushBefore := func(ln line) {
		if beforeN == 0 {
			return
		}
		if len(beforeBuf) < beforeN {
			beforeBuf = append(beforeBuf, ln)
			return
		}
		copy(beforeBuf, beforeBuf[1:])
		beforeBuf[len(beforeBuf)-1] = ln
	}

	lineNo := 0
	for sc.Scan() {
		lineNo++
		txt := sc.Text()
		ln := line{no: lineNo, text: txt}

		m := matchFn(txt)
		if o.invert {
			m = !m
		}

		if m {

			for _, b := range beforeBuf {
				if b.no > lastPrintedNo {
					printLine(b, false)
				}
			}

			if ln.no > lastPrintedNo {
				printLine(ln, true)
			}

			if afterN > afterRemain {
				afterRemain = afterN
			}
		} else if afterRemain > 0 {

			if ln.no > lastPrintedNo {
				printLine(ln, false)
			}
			afterRemain--
		}

		pushBefore(ln)
	}

	return sc.Err()
}

func newScanner(r io.Reader) *bufio.Scanner {
	sc := bufio.NewScanner(r)

	const maxLine = 4 << 20
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxLine)
	return sc
}
