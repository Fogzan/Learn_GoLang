package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func firstGetData(ctx context.Context, ch chan int) {
	for {
		select {
		case <-time.After(200 * time.Millisecond):
		case <-ctx.Done():
			return
		}
		select {
		case ch <- rand.Intn(99):
		case <-ctx.Done():
			return
		}
	}

}

func secondGetData(ctx context.Context, ch chan int) {
	for {
		select {
		case <-time.After(300 * time.Millisecond):
		case <-ctx.Done():
			return
		}
		select {
		case ch <- 100 + rand.Intn(99):
		case <-ctx.Done():
			return
		}
	}

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ch := make(chan int)
	go firstGetData(ctx, ch)
	go secondGetData(ctx, ch)

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-ch:
			fmt.Println(result)
		}
	}
}
