package main

import "fmt"

func main() {
	str := "главрыба"
	str = reverseStr(str)
	fmt.Println(str)

}

func reverseStr(str string) string {
	symbs := []rune(str)
	for i, j := 0, len(symbs)-1; i < j; i++ {
		symbs[i], symbs[j] = symbs[j], symbs[i]
		j--
	}
	return string(symbs)
}
