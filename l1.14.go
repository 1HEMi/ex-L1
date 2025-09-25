package main

import (
	"fmt"
)

func main() {
	DetectType(1)

}

func DetectType(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Println("This is int:", val)
	case string:
		fmt.Println("This is string:", val)
	case bool:
		fmt.Println("This is bool:", val)
	case chan int:
		fmt.Println("This is int chan:", val)
	case chan string:
		fmt.Println("This is string chan:", val)
	case chan struct{}:
		fmt.Println("This is struct chan:", val)

	default:
		fmt.Println("unrecognized type:", val)
	}
}
