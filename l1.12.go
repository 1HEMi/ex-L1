package main

import "fmt"

func main() {
	arr := []string{"cat", "cat", "dog", "cat", "tree"}
	res := Unique(arr)
	fmt.Println(res)

}

func Unique(arr []string) []string {
	unique := make(map[string]bool)
	res := []string{}
	for _, v := range arr {
		unique[v] = true

	}
	for k := range unique {
		res = append(res, k)
	}
	return res
}
