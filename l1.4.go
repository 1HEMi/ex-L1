package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func worker(id int, ctx context.Context, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Printf("worker %d: channel closed\n", id)
				return
			}
			fmt.Printf("worker %d got: %d\n", id, v)
		case <-ctx.Done():
			fmt.Printf("worker %d: ctx canceled\n", id)
			return
		}
	}
}

func main() {
	numWorkers := flag.Int("n", 2, "number of workers")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(*numWorkers)
	for i := 1; i <= *numWorkers; i++ {
		go worker(i, ctx, ch, &wg)
	}

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			close(ch)
			wg.Wait()
			fmt.Println("shutdown complete")
			return
		case <-time.After(400 * time.Millisecond):
			ch <- i
		}
	}
}
