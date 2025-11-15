package main

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidFormat = errors.New("invalid string format")
)

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	var (
		b        strings.Builder
		prevRune rune
		hasPrev  bool
		escaped  bool
		i        int
	)
	b.Grow(len(s))

	flushPrevOnce := func() {
		if hasPrev {
			b.WriteRune(prevRune)
			hasPrev = false
		}
	}
	emitPrevTimes := func(n int) error {
		if !hasPrev {
			return ErrInvalidFormat
		}
		if n < 0 {
			return ErrInvalidFormat
		}
		for j := 0; j < n; j++ {
			b.WriteRune(prevRune)
		}
		hasPrev = false
		return nil
	}

	for i < len(s) {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {

			return "", ErrInvalidFormat
		}

		if escaped {
			flushPrevOnce()
			prevRune = r
			hasPrev = true
			escaped = false
			i += size
			continue
		}

		if r == '\\' {
			escaped = true
			i += size
			continue
		}

		if unicode.IsDigit(r) {

			if !hasPrev {
				return "", ErrInvalidFormat
			}

			count := int(r - '0')
			i += size
			for i < len(s) {
				r2, sz2 := utf8.DecodeRuneInString(s[i:])
				if !unicode.IsDigit(r2) {
					break
				}
				count = count*10 + int(r2-'0')
				i += sz2
			}
			if err := emitPrevTimes(count); err != nil {
				return "", err
			}
			continue
		}

		flushPrevOnce()
		prevRune = r
		hasPrev = true
		i += size
	}

	if escaped {
		return "", ErrInvalidFormat
	}
	flushPrevOnce()
	return b.String(), nil
}
