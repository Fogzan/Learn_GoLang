package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func processData(ctx context.Context, val int) (int, error) {
	ch := make(chan struct{})

	go func() {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		close(ch)
	}()

	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	case <-ch:
		return val * 2, nil
	}

}

func main() {
	in := make(chan int)
	out := make(chan int)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		// time.Sleep(10 * time.Second)
		for i := range 100 {
			select {
			case in <- i:
			case <-ctx.Done():
				return
			}
		}
		close(in)
	}()

	now := time.Now()

	processParallel(ctx, in, out, 5)

	for val := range out {
		fmt.Println(val)
	}

	fmt.Println(time.Since(now))
}

func processParallel(ctx context.Context, in <-chan int, out chan<- int, numWorkers int) {
	var wgFirst sync.WaitGroup

	wgFirst.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wgFirst.Done()

			for {
				select {
				case val, ok := <-in:
					if !ok {
						return
					}
					result, err := processData(ctx, val)

					if err != nil {
						return
					} else {
						out <- result
					}

				case <-ctx.Done():
					return
				}

			}
		}()
	}
	go func() {
		wgFirst.Wait()
		close(out)
	}()
}
