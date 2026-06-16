package main

import (
	"context"
	"fmt"
	"time"
)

func Generate(ctx context.Context, nums ...int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for _, i := range nums {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}

		}
	}()

	return ch
}

func Square(ctx context.Context, in <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
	mainFor:
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-in:
				if !ok {
					break mainFor
				}
				select {
				case <-ctx.Done():
					return
				case ch <- i * i:
				}
			}

		}
	}()

	return ch
}

func Print(ctx context.Context, in <-chan int) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)
	mainFor:
		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-in:
				if !ok {
					break mainFor
				}
				time.Sleep(500 * time.Millisecond)
				fmt.Println(i)
			}
		}
	}()

	return ch
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	timeStart := time.Now()
	defer func() {
		fmt.Printf("Время выполнения паттерна: %v\n", time.Since(timeStart))
	}()

	array := []int{5, 6, 7, 8, 4, 3, 2, 5, 6, 76}

	waitChan := Print(ctx, Square(ctx, Generate(ctx, array...)))
	select {
	case <-ctx.Done():
		return
	case <-waitChan:
	}

}
