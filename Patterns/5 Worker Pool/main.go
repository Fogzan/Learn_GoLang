package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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
				time.Sleep(500 * time.Millisecond)
				fmt.Printf("Worker: %v обработал число: %v\n", id, i)
				select {
				case <-ctx.Done():
					return
				case ch <- i * 2:
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

func FanIn(ctx context.Context, channels ...<-chan int) chan int {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, elem := range channels {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-elem:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case ch <- data:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func WorkerPool(ctx context.Context, jobs <-chan int, workers int) <-chan int {
	ch := FanIn(ctx, FanOut(ctx, jobs, workers)...)
	return ch
}

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	// Завершает работу при сигнале
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()
		fmt.Println("SAVE DATA")
	}()

	timeStart := time.Now()
	defer func() {
		fmt.Printf("Время выполнения паттерна: %v\n", time.Since(timeStart))
	}()

	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := range 100 {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}

		}
	}()

	resultChan := WorkerPool(ctx, ch, 5)

mainFor:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ПО КОНТЕКСТУ")
			break mainFor
		case i, ok := <-resultChan:
			if !ok {
				break mainFor
			}
			fmt.Printf("Result: %v\n", i)
		}
	}

}
