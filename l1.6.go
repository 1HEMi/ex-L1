package main

import (
	//"context"
	"fmt"
	//"runtime"
	"time"
)

// 	// Выход из горутины по условию
// func main() {
// 	go func() {
// 		for i := 0; i < 10; i++ {
// 			fmt.Println("work", i)

// 		}
// 		fmt.Println("goroutine done by condition")
// 	}()

// 	time.Sleep(3 * time.Second)
// }

// Выход из горутины по каналу
// func main() {

// 	done := make(chan struct{})

// 	go func() {
// 		for {
// 			select {
// 			case <-done:
// 				fmt.Println("goroutine stopped by channel")
// 				return
// 			default:
// 				fmt.Println("working")
// 				time.Sleep(500 * time.Millisecond)
// 			}
// 		}
// 	}()

// 	time.Sleep(2 * time.Second)
// 	done <- struct{}{}
// 	time.Sleep(time.Second)
// }

// Выход из горутины с помощью context
// func main() {
// 	ctx, cancel := context.WithCancel(context.Background())

// 	go func(ctx context.Context) {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				fmt.Println("goroutine stopped by context")
// 				return
// 			default:
// 				fmt.Println("working")
// 				time.Sleep(500 * time.Millisecond)
// 			}
// 		}
// 	}(ctx)

// 	time.Sleep(2 * time.Second)
// 	cancel()
// 	time.Sleep(time.Second)
// }

// Выход из горутины с помощью runtime.Goexit
// func main() {
// 	go func() {
// 		fmt.Println("goroutine starting")
// 		time.Sleep(time.Second)
// 		fmt.Println("goroutine calling Goexit()")
// 		runtime.Goexit()

// 	}()

// 	time.Sleep(2 * time.Second)
// 	fmt.Println("main done")
// }

// Выход из горутины с помощью закрытия канала
func main() {
	ch := make(chan int)

	go func() {
		for v := range ch {
			fmt.Println("got:", v)
		}
		fmt.Println("goroutine stopped because channel closed")
	}()

	for i := 0; i < 3; i++ {
		ch <- i
		time.Sleep(500 * time.Millisecond)
	}
	close(ch)
	time.Sleep(time.Second)
}
