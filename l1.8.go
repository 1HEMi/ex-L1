package main

import (
	"fmt"
)

func main() {
	var num int64 = 5
	var i uint = 0

	fmt.Printf("Before:= %d ", num)
	num = SetBit(num, i, 0)
	fmt.Printf("After:= %d", num)

}

func SetBit(num int64, i uint, bit int) int64 {
	if bit == 1 {
		num |= 1 << i
	} else {
		num &^= 1 << i
	}
	return num
}
