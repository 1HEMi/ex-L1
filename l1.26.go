package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(IsUnique("abcd"))
	fmt.Println(IsUnique("aAbcd"))
	fmt.Println(IsUnique("aabcd"))
}

func IsUnique(s string) bool {
	s = strings.ToLower(s)
	unique := make(map[rune]bool)

	for _, ch := range s {
		if unique[ch] {
			return false
		}
		unique[ch] = true
	}
	return true
}
