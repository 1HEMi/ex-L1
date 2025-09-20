package main

import (
	"fmt"
)

func main() {
	nums := []int{1, 2, 3, 4, 5}
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func(out chan<- int) {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}(ch1)

	go func(in <-chan int, out chan<- int) {
		defer close(out)
		for x := range in {
			out <- x * 2
		}
	}(ch1, ch2)

	for v := range ch2 {
		fmt.Println(v)
	}
}
