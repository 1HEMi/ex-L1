package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type options struct {
	col       int  // -k
	numeric   bool // -n
	reverse   bool // -r
	unique    bool // -u
	month     bool // -M
	trimRight bool // -b
	checkOnly bool // -c
	human     bool // -h
}

func main() {
	opts := parseFlags()

	var in io.Reader = os.Stdin
	args := flag.Args()
	if len(args) > 0 && args[0] != "-" {
		f, err := os.Open(args[0])
		if err != nil {
			fail(err)
		}
		defer f.Close()
		in = f
	}

	lines, err := readAllLines(in)
	if err != nil {
		fail(err)
	}

	less := makeLess(lines, opts)

	if opts.checkOnly {
		if dis, i := firstDisorder(lines, less); dis {
			fmt.Fprintf(os.Stderr, "data not sorted at line %d\n", i+2)
			os.Exit(1)
		}
		return
	}

	sort.SliceStable(lines, func(i, j int) bool { return less(i, j) })

	if opts.unique && len(lines) > 0 {
		lines = dedupe(lines)
	}

	for _, s := range lines {
		fmt.Println(s)
	}
}

func parseFlags() options {
	var col int
	flag.IntVar(&col, "k", 0, "column number (1-based) to sort by; 0 means whole line (tab-delimited)")
	var numeric = flag.Bool("n", false, "compare by numeric value")
	var reverse = flag.Bool("r", false, "reverse the result of comparisons")
	var unique = flag.Bool("u", false, "output only the first of an equal run")
	var month = flag.Bool("M", false, "compare by month name (Jan..Dec)")
	var trim = flag.Bool("b", false, "ignore trailing blanks in key")
	var check = flag.Bool("c", false, "check whether input is sorted; do not sort")
	var human = flag.Bool("h", false, "sort by numeric value, taking into account suffixes (for example, K = kilobytes, M = megabytes - human-readable sizes")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [file|-]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if col < 0 {
		fail(errors.New("column -k must be >= 0"))
	}
	return options{
		col:       col,
		numeric:   *numeric,
		reverse:   *reverse,
		unique:    *unique,
		month:     *month,
		trimRight: *trim,
		checkOnly: *check,
		human:     *human,
	}
}

func readAllLines(r io.Reader) ([]string, error) {
	sc := bufio.NewScanner(r)

	const maxCap = 4 << 20
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxCap)

	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

func makeLess(lines []string, opts options) func(i, j int) bool {
	key := func(s string) string {
		if opts.col <= 0 {
			return s
		}
		cols := strings.Split(s, "\t")
		if opts.col-1 >= len(cols) {
			return ""
		}
		k := cols[opts.col-1]
		if opts.trimRight {
			k = strings.TrimRight(k, " \t")
		}
		return k
	}

	var cmp func(a, b string) int

	switch {
	case opts.human:
		cmp = func(a, b string) int {
			av, aok := parseHuman(key(a))
			bv, bok := parseHuman(key(b))

			if !aok || !bok {
				return strings.Compare(key(a), key(b))
			}
			if av < bv {
				return -1
			}
			if av > bv {
				return 1
			}
			return 0
		}
	case opts.numeric:
		cmp = func(a, b string) int {
			af, aok := parseFloat(key(a))
			bf, bok := parseFloat(key(b))
			if !aok || !bok {
				return strings.Compare(key(a), key(b))
			}
			if af < bf {
				return -1
			}
			if af > bf {
				return 1
			}
			return 0
		}
	case opts.month:
		cmp = func(a, b string) int {
			am, aok := monthIndex(key(a))
			bm, bok := monthIndex(key(b))
			if !aok || !bok {
				return strings.Compare(key(a), key(b))
			}
			if am < bm {
				return -1
			}
			if am > bm {
				return 1
			}
			return 0
		}
	default:
		cmp = func(a, b string) int {
			return strings.Compare(key(a), key(b))
		}
	}

	return func(i, j int) bool {
		res := cmp(lines[i], lines[j])
		if opts.reverse {
			return res > 0
		}
		return res < 0
	}
}

func firstDisorder(lines []string, less func(i, j int) bool) (bool, int) {
	for i := 0; i+1 < len(lines); i++ {
		if less(i+1, i) {
			return true, i
		}
	}
	return false, -1
}

func dedupe(sorted []string) []string {
	if len(sorted) == 0 {
		return sorted
	}
	out := sorted[:1]
	last := sorted[0]
	for _, s := range sorted[1:] {
		if s != last {
			out = append(out, s)
			last = s
		}
	}
	return out
}

func parseFloat(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil
}

func parseHuman(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	i := 0
	for i < len(s) && (unicode.IsDigit(rune(s[i])) || s[i] == '.' || s[i] == '+' || s[i] == '-') {
		i++
	}
	numPart := s[:i]
	suf := strings.TrimSpace(strings.TrimSuffix(strings.ToUpper(s[i:]), "B"))

	val, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, false
	}
	mult := 1.0
	switch suf {
	case "":
		mult = 1
	case "K":
		mult = 1024
	case "M":
		mult = 1024 * 1024
	case "G":
		mult = 1024 * 1024 * 1024
	case "T":
		mult = 1024 * 1024 * 1024 * 1024
	case "P":
		mult = 1024 * 1024 * 1024 * 1024 * 1024
	case "E":
		mult = 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return 0, false
	}
	return val * mult, true
}

func monthIndex(s string) (int, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case "jan", "january":
		return 1, true
	case "feb", "february":
		return 2, true
	case "mar", "march":
		return 3, true
	case "apr", "april":
		return 4, true
	case "may":
		return 5, true
	case "jun", "june":
		return 6, true
	case "jul", "july":
		return 7, true
	case "aug", "august":
		return 8, true
	case "sep", "september":
		return 9, true
	case "oct", "october":
		return 10, true
	case "nov", "november":
		return 11, true
	case "dec", "december":
		return 12, true
	default:
		return 0, false
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
