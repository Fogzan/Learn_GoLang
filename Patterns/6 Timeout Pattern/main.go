package main

import (
	"context"
	"fmt"
	"time"
)

func worker(input int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		time.Sleep(1 * time.Second)
		ch <- input * 5
	}()

	return ch
}

func WithTimeout(ctx context.Context, input <-chan int, timeout time.Duration) <-chan int {
	ch := make(chan int)

	timeoutCh := time.After(timeout)

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			case <-timeoutCh:
				return
			case data, ok := <-input:
				if !ok {
					return
				}
				chanWorker := worker(data)
				select {
				case <-ctx.Done():
					return
				case <-timeoutCh:
					return
				case result, okResult := <-chanWorker:
					if !okResult {
						return
					}
					select {
					case <-ctx.Done():
						return
					case <-timeoutCh:
						return
					case ch <- result:
					}
				}
			}
		}

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

	ch := make(chan int)
	go func() {
		defer close(ch)

		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()

	chOut := WithTimeout(ctx, ch, 4*time.Second)
	// chOut := WithTimeout(ctx, ch, 3*time.Millisecond)
mainFor:
	for {
		select {
		case <-ctx.Done():
			return
		case result, ok := <-chOut:
			if !ok {
				break mainFor
			}
			fmt.Printf("Результат: %v\n", result)
		}
	}

}
