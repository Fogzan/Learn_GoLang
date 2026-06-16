package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Worker(ctx context.Context, id int, in <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case <-time.After(1 * time.Second):
					fmt.Printf("Worker: %v обработал число: %v\n", id, i)
					select {
					case <-ctx.Done():
						return
					case ch <- i * 2:
					}
				}

			}
		}
	}()

	return ch
}

func FanOut(ctx context.Context, in <-chan int, workers int) []<-chan int {
	chans := make([]<-chan int, workers)
	for i := range workers {
		chans[i] = Worker(ctx, i, in)
	}
	return chans
}

func main() {
	timeStart := time.Now()
	defer func() {
		fmt.Printf("Время выполнения паттерна: %v\n", time.Since(timeStart))
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := range 10 {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()

	workers := FanOut(ctx, ch, 2)

	wg := sync.WaitGroup{}

	for _, elem := range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case i, ok := <-elem:
					if !ok {
						return
					}
					fmt.Println(i)
				}
			}
		}()
	}
	wg.Wait()

}
