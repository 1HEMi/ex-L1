package main

import "fmt"

func main() {
	arr := []int{5, 3, 8, 4, 2, 7, 1, 10, 6}
	sortedArr := quickSort(arr)
	fmt.Println(sortedArr)

}

func quickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	pivot := arr[len(arr)/2]
	var left, right []int
	for i, v := range arr {
		if i == len(arr)/2 {
			continue
		}
		if v <= pivot {
			left = append(left, v)
		} else {
			right = append(right, v)
		}

	}
	return append(append(quickSort(left), pivot), quickSort(right)...)

}
