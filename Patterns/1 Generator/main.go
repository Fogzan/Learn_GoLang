package main

import (
	"context"
	"fmt"
	"time"
)

func NumbersGenerator(ctx context.Context, n int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := range n {
			select {
			case <-ctx.Done():
				return
			case ch <- i + 1:
			}

		}
	}()

	return ch
}

func FilterEven(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

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
				case out <- i + 10:

				}
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch := NumbersGenerator(ctx, 10)
	chSecond := FilterEven(ctx, ch)
	for i := range chSecond {
		fmt.Println(i)
	}
}
