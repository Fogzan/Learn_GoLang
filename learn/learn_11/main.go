package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	fmt.Println("Приложение завершило работу!")
}
