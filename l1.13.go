package main

import "fmt"

func main() {

	a := 2
	b := 10

	fmt.Printf("a = %d, b = %d\n", a, b)
	a, b = Swap(a, b)
	fmt.Printf("a = %d, b = %d\n", a, b)

	c := 5
	d := 17
	fmt.Println("------------------")
	fmt.Printf("c = %d, d = %d\n", c, d)
	c, d = XORSwap(c, d)

	fmt.Printf("c = %d, d = %d", c, d)

}

func Swap(a, b int) (int, int) {
	a = a + b
	b = a - b
	a = a - b
	return a, b
}

func XORSwap(a, b int) (int, int) {
	a = a ^ b
	b = a ^ b
	a = a ^ b

	return a, b
}
