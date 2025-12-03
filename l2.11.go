package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	result := Anagram(words)
	fmt.Println(result)
}

func Anagram(words []string) map[string][]string {
	anagrams := make(map[string][]string)
	firstWord := make(map[string]string)
	for _, val := range words {
		lowerVal := strings.ToLower(val)
		runes := []rune(lowerVal)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		canonical := string(runes)
		if _, ok := firstWord[canonical]; !ok {
			firstWord[canonical] = lowerVal
		}
		anagrams[canonical] = append(anagrams[canonical], lowerVal)

	}
	result := make(map[string][]string)
	for canonical, group := range anagrams {
		if len(group) < 2 {
			continue
		}
		sort.Strings(group)
		key := firstWord[canonical]
		result[key] = group

	}

	return result
}
