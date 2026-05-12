package main

import (
	"fmt"
	"time"
)

func writer() <-chan int {
	chFirst := make(chan int)

	go func() {
		for i := range 10 {
			chFirst <- i + 1
		}
		close(chFirst)
	}()

	return chFirst
}

func doubler(chFirst <-chan int) <-chan int {
	chSecond := make(chan int)

	go func() {
		for elem := range chFirst {
			time.Sleep(500 * time.Millisecond)
			chSecond <- elem * 2
		}
		close(chSecond)
	}()

	return chSecond
}

func reader(chSecond <-chan int) {
	for elem := range chSecond {
		fmt.Println(elem)
	}
}

// writer - генерирует числи от 1 до 10
// doubler - умножает числа на 2, задержка 500мс
// reader - выводит числа

func main() {
	reader(doubler(writer()))
}
