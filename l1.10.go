package main

import "fmt"

func main() {
	weather := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	arr1 := []float64{}
	arr2 := []float64{}
	arr3 := []float64{}
	arr4 := []float64{}

	for _, v := range weather {
		switch {
		case v >= -30 && v < -20:
			arr1 = append(arr1, v)
		case v >= 10 && v < 20:
			arr2 = append(arr2, v)
		case v >= 20 && v < 30:
			arr3 = append(arr3, v)
		case v >= 30 && v < 40:
			arr4 = append(arr4, v)
		}

	}

	fmt.Println("-20:", arr1)
	fmt.Println("10:", arr2)
	fmt.Println("20:", arr3)
	fmt.Println("30:", arr4)

}
