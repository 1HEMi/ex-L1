package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	target := binarySearch(arr, 11)
	fmt.Println(target)

}

func binarySearch(arr []int, num int) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] == num {
			return mid
		}
		if arr[mid] < num {
			low = mid + 1
		} else {
			high = mid - 1
		}

	}
	return -1
}
