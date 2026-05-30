package main

import (
	"context"
	"fmt"
	"sync"
	"tee/teechan"
	"time"
)

func generate(ctx context.Context) chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := range 5 {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()

	return ch
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(3)

	// chanArray := teechan(ctx, generate(ctx), 3)
	teechanElem := teechan.New(3, teechan.Fast)
	chanArray := teechanElem.Execute(ctx, generate(ctx))

	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-chanArray[0]:
				if !ok {
					return
				}
				fmt.Println("Обработчик: ", val)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-chanArray[1]:
				if !ok {
					return
				}
				fmt.Println("Логгер: ", val)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-chanArray[2]:
				if !ok {
					return
				}
				time.Sleep(time.Second)
				fmt.Println("БД: ", val)
			}
		}
	}()

	wg.Wait()
}
