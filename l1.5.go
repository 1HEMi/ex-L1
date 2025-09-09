package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		for i := 0; ; i++ {
			ch <- i

		}
	}()

	N := flag.Int("N", 5, "program running time")
	flag.Parse()
	timeout := time.After(time.Duration(*N) * time.Second)

	for {
		select {
		case v := <-ch:
			fmt.Println("got:", v)
		case <-timeout:
			fmt.Println("time is up, stopping program...")
			close(ch)
			return

		}
	}
}
