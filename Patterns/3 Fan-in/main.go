package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

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

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	timeStart := time.Now()
	defer func() {
		fmt.Printf("Время выполнения паттерна: %v\n", time.Since(timeStart))
	}()

	countChan := 2
	channels := make([]chan int, countChan)
	for i := range countChan {
		channels[i] = make(chan int)
	}

	go func() {
		defer close(channels[0])
		for i := range 100 {
			select {
			case <-ctx.Done():
				return
			case channels[0] <- i + 100:
				time.Sleep(500 * time.Millisecond)
			}

		}
	}()

	go func() {
		defer close(channels[1])
		for i := range 150 {
			select {
			case <-ctx.Done():
				return
			case channels[1] <- i:
			}
		}

	}()

	channelsToFunc := make([]<-chan int, countChan)
	for i := range countChan {
		channelsToFunc[i] = channels[i]
	}

	resultChan := FanIn(ctx, channelsToFunc...)

mainfor:
	for {
		select {
		case <-ctx.Done():
			return
		case i, ok := <-resultChan:
			if !ok {
				break mainfor
			}
			fmt.Printf("%v\n", i)
		}
	}

}
