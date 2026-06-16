package main

import (
	"context"
	"fmt"
	"time"
)

func Merge(ctx context.Context, a, b <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for a != nil || b != nil {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v

			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	timeStart := time.Now()
	defer func() {
		fmt.Printf("Время выполнения: %v\n", time.Since(timeStart))
	}()

	a := make(chan int)
	b := make(chan int)

	go func() {
		defer close(a)

		for i := 0; i < 5; i++ {
			a <- i
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		defer close(b)

		for i := 100; i < 105; i++ {
			b <- i
			time.Sleep(300 * time.Millisecond)
		}
	}()

	for v := range Merge(ctx, a, b) {
		fmt.Println(v)
	}
}
