package main

import (
	"fmt"
	"time"
)

func MySleepTimer(d time.Duration) {
	timer := time.NewTimer(d)
	<-timer.C
}

func main() {
	fmt.Println("Start")
	MySleepTimer(2 * time.Second)
	fmt.Println("End")
}
