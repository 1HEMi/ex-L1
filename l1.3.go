package main

import (
	"flag"
	"fmt"
	"sync"
)

func worker(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("worker %d got: %d\n", id, v)
	}
}

func main() {
	numWorkers := flag.Int("n", 2, "number of workers")
	flag.Parse()

	var wg sync.WaitGroup
	ch := make(chan int)

	for i := 1; i <= *numWorkers; i++ {
		wg.Add(1)
		go worker(i, ch, &wg)
	}

	for i := 0; ; i++ {
		ch <- i
	}

}
