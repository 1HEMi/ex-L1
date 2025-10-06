package main

import "fmt"

func RemoveAtIndex(slice []int, i int) []int {
	if i < 0 || i >= len(slice) {
		return slice
	}
	copy(slice[i:], slice[i+1:])
	slice = slice[:len(slice)-1]

	return slice
}

func main() {
	arr := []int{10, 20, 30, 40, 50}
	fmt.Println(arr)

	arr = RemoveAtIndex(arr, 2)

	fmt.Println(arr)
}
