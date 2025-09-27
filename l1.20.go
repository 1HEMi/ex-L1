package main

import (
	"fmt"
	"strings"
)

func main() {
	sentence := "snow dog sun"
	res := ReverseSentence(sentence)
	fmt.Println(res)

}

func ReverseSentence(sentence string) string {

	words := strings.Fields(sentence)

	for i, j := 0, len(words)-1; i < j; i++ {
		words[i], words[j] = words[j], words[i]
		j--
	}

	res := strings.Join(words, " ")
	return res
}
