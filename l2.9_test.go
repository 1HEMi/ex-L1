package main

import (
	"testing"
)

func TestUnpack_Basic(t *testing.T) {
	tests := []struct {
		in    string
		out   string
		valid bool
	}{
		{"a4bc2d5e", "aaaabccddddde", true},
		{"abcd", "abcd", true},
		{"", "", true},
		{"45", "", false},    // начинается с цифры
		{`\4\5`, "45", true}, // экранированные цифры — литералы
		{`qwe\4\5`, "qwe45", true},
		{`qwe\45`, "qwe44444", true},  // \4 — литерал '4', затем повторить 5 раз
		{`a0b`, "b", true},            // ноль — просто не добавляем 'a'
		{`й2`, "йй", true},            // юникод-руны
		{`a12`, "aaaaaaaaaaaa", true}, // многозначный счётчик
		{`a\`, "", false},             // висящий слеш
	}

	for _, tc := range tests {
		got, err := Unpack(tc.in)
		if tc.valid && err != nil {
			t.Fatalf("Unpack(%q) unexpected error: %v", tc.in, err)
		}
		if !tc.valid && err == nil {
			t.Fatalf("Unpack(%q) expected error, got nil", tc.in)
		}
		if tc.valid && got != tc.out {
			t.Fatalf("Unpack(%q) = %q, want %q", tc.in, got, tc.out)
		}
	}
}
