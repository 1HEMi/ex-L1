package main

import "fmt"

func main() {
	a := []int{1, 2, 3}
	b := []int{2, 3, 4}
	res := Intersection(a, b)
	fmt.Println("intersection:", res)

}

func Intersection(arr1 []int, arr2 []int) []int {
	intersection := []int{}
	for _, v := range arr1 {
		for _, v2 := range arr2 {
			if v == v2 {
				intersection = append(intersection, v)
			}
		}
	}
	return intersection

}
