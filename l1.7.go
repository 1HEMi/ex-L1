package main

import (
	"sync"
)

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	m := make(map[int]string)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			writeToMap(m, val, "number", &mu)
		}(i)
	}

	wg.Wait()

}

func writeToMap(m map[int]string, key int, value string, mu *sync.Mutex) {

	mu.Lock()
	defer mu.Unlock()
	m[key] = value

}
