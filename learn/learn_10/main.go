package main

import (
	"fmt"
	"sync"
)

func gooodFuncExample(in <-chan int) (<-chan int, <-chan int) {
	firstOut := make(chan int)
	secondOut := make(chan int)

	go func() {
		firstWg := sync.WaitGroup{}
		secondWg := sync.WaitGroup{}
		for result := range in {
			firstWg.Add(1)
			secondWg.Add(1)
			go func() {
				defer firstWg.Done()
				firstOut <- result
			}()
			go func() {
				defer secondWg.Done()
				secondOut <- result
			}()
		}
		go func() {
			firstWg.Wait()
			close(firstOut)
		}()
		go func() {
			secondWg.Wait()
			close(secondOut)
		}()
	}()

	return firstOut, secondOut
}

func main() {

	in := make(chan int)
	go func() {
		for i := range 10 {
			in <- i
		}
		close(in)
	}()

	firstOut, secondOut := gooodFuncExample(in)
	for i := range firstOut {
		fmt.Printf("Первый канал, значение: %v\n", i)
	}
	for i := range secondOut {
		fmt.Printf("Второй канал, значение: %v\n", i)
	}
}
