package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randomWait() int {
	workSeconds := rand.Intn(5 + 1)
	time.Sleep(time.Duration(workSeconds) * time.Second)
	return workSeconds
}

func main() {
	ch := make(chan int)

	resultSeconds := 0
	start := time.Now()

	for range 100 {
		go func() {
			ch <- randomWait()
		}()
	}

	for range 100 {
		resultSeconds += <-ch
	}

	fmt.Printf("All time: \t%vs\n", resultSeconds)
	fmt.Printf("Main time: \t%v\n", time.Since(start))
}
