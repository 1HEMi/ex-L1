package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	val int
	mx  sync.Mutex
}

func main() {
	var wg sync.WaitGroup
	wg.Add(10)
	var counter Counter
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Println("result:", counter.val)
}

func (c *Counter) Increment() {
	c.mx.Lock()
	c.val++
	c.mx.Unlock()

}
