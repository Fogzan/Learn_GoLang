package main

import (
	"fmt"
	"sync"
	"time"
)

func Worker(id int, in <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for data := range in {
			time.Sleep(1 * time.Second)
			fmt.Printf("Worker: %v обработал число: %v\n", id, data)
			ch <- data * 2
		}

	}()

	return ch
}

func FanOut(in <-chan int, workers int) []<-chan int {
	chans := make([]<-chan int, workers)
	for i := range workers {
		chans[i] = Worker(i, in)
	}
	return chans
}

func main() {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := range 100 {
			ch <- i
		}
	}()

	workerPool := FanOut(ch, 2)

	wg := sync.WaitGroup{}

	for _, elem := range workerPool {
		wg.Add(1)
		go func() {
			for i := range elem {
				fmt.Println(i)
			}
			wg.Done()
		}()
	}
	wg.Wait()

}
