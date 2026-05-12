package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

func randomTimeWork() {
	time.Sleep(time.Duration(rand.IntN(100)) * time.Second)
}

func predictableTimeWork(mainFunc func()) error {
	ch := make(chan struct{})

	go func() {
		mainFunc()
		close(ch)
	}()

	select {
	case <-ch:
		fmt.Println("Успешно")
		return nil
	case <-time.After(3 * time.Second):
		fmt.Println("Ошибка")
		return errors.New("Error. Function worker big 3 srconds!")
	}
}

func main() {
	predictableTimeWork(randomTimeWork)
}
