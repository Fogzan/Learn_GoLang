package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func WorkerPool(ctx context.Context, jobs <-chan int, workers int) <-chan int {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := range workers {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-jobs:
					if !ok {
						return
					}
					time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
					fmt.Printf("Обработчик %v обработал %v\n", i, data)
					select {
					case <-ctx.Done():
						return
					case ch <- data * 2:
					}
				}
			}

		}()
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	timeStart := time.Now()
	defer func() {
		fmt.Printf("Время выполнения: %v\n", time.Since(timeStart))
	}()

	jobs := make(chan int)

	go func() {
		defer close(jobs)
		for i := range 100 {
			select {
			case <-ctx.Done():
				return
			case jobs <- i:
			}
		}
	}()

	result := WorkerPool(ctx, jobs, 5)

mainfor:
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-result:
			if !ok {
				break mainfor
			}
			fmt.Printf("Результат обработки: %v\n", val)
		}
	}
}
