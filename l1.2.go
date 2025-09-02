package main

import (
	"fmt"
	"sync"
)

func main() {
	var digits = [5]int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}
	wg.Add(len(digits))

	for _, v := range digits {
		go func(val int) {
			defer wg.Done()
			fmt.Println(val * val)
		}(v)
	}
	wg.Wait()

}
