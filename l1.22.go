package main

import (
	"fmt"
	"math/big"
)

func main() {
	res := Calculator("912331231213", "32353343434", "/")
	fmt.Println(res)
}

func Calculator(aStr, bStr, op string) *big.Int {
	a := new(big.Int)
	b := new(big.Int)
	a.SetString(aStr, 10)
	b.SetString(bStr, 10)
	result := new(big.Int)

	switch op {
	case "+":
		result.Add(a, b)
	case "-":
		result.Sub(a, b)
	case "*":
		result.Mul(a, b)
	case "/":
		result.Div(a, b)
	default:
		fmt.Println("Invalid operation")
		return nil
	}

	return result
}
