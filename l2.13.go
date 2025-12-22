// Program cut reads lines from STDIN, splits by delimiter and prints selected fields.
// Flags:
//
//	-f "fields"   fields to select: comma-separated numbers and ranges (e.g. "1,3-5,8")
//	-d "delimiter"    delimiter (single character), default is tab ('\t')
//	-s            only print lines containing the delimiter (skip lines without delimiter)
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type options struct {
	fieldsSpec string
	delimiter  string
	separated  bool
}

func main() {
	var o options
	flag.StringVar(&o.fieldsSpec, "f", "", "fields to select (e.g. \"1,3-5\")")
	flag.StringVar(&o.delimiter, "d", "\t", "delimiter (single character), default is TAB")
	flag.BoolVar(&o.separated, "s", false, "only print lines containing the delimiter")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -f <fields> [-d <delim>] [-s]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if strings.TrimSpace(o.fieldsSpec) == "" {
		dieUsage("missing -f fields specification")
	}
	if len([]rune(o.delimiter)) != 1 {
		dieUsage("delimiter must be a single character")
	}
	delim := o.delimiter

	fields, err := parseFields(o.fieldsSpec)
	if err != nil {
		dieErr(err, 2)
	}

	if err := cut(os.Stdin, os.Stdout, delim, fields, o.separated); err != nil {
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

func cut(r io.Reader, w io.Writer, delim string, fields []int, separated bool) error {
	sc := bufio.NewScanner(r)
	const maxLine = 4 << 20
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, maxLine)

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for sc.Scan() {
		line := sc.Text()

		if separated && !strings.Contains(line, delim) {
			continue
		}

		parts := strings.Split(line, delim)

		out := make([]string, 0, len(fields))
		for _, f := range fields {
			idx := f - 1
			if idx >= 0 && idx < len(parts) {
				out = append(out, parts[idx])
			}
		}

		_, err := fmt.Fprintln(bw, strings.Join(out, delim))
		if err != nil {
			return err
		}
	}
	return sc.Err()
}

func parseFields(spec string) ([]int, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, errors.New("empty fields spec")
	}

	set := make(map[int]struct{})
	chunks := strings.Split(spec, ",")
	for _, c := range chunks {
		c = strings.TrimSpace(c)
		if c == "" {
			return nil, errors.New("invalid fields spec: empty segment")
		}

		if strings.Contains(c, "-") {
			a, b, ok := strings.Cut(c, "-")
			if !ok {
				return nil, errors.New("invalid range")
			}
			start, err := parsePosInt(strings.TrimSpace(a))
			if err != nil {
				return nil, err
			}
			end, err := parsePosInt(strings.TrimSpace(b))
			if err != nil {
				return nil, err
			}
			if start > end {
				return nil, errors.New("invalid range: start > end")
			}
			for i := start; i <= end; i++ {
				set[i] = struct{}{}
			}
			continue
		}

		n, err := parsePosInt(c)
		if err != nil {
			return nil, err
		}
		set[n] = struct{}{}
	}

	fields := make([]int, 0, len(set))
	for n := range set {
		fields = append(fields, n)
	}
	sort.Ints(fields)
	return fields, nil
}

func parsePosInt(s string) (int, error) {
	if s == "" {
		return 0, errors.New("invalid number: empty")
	}
	n := 0
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < '0' || ch > '9' {
			return 0, errors.New("invalid number: " + s)
		}
		n = n*10 + int(ch-'0')
	}
	if n <= 0 {
		return 0, errors.New("field numbers must be >= 1")
	}
	return n, nil
}
