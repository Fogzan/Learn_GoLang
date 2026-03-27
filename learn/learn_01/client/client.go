package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				break
			}
			if string(buff[0:n]) != "true" {
				fmt.Print(string(buff[0:n]))
			} else {
				wg.Done()
				conn.Close()
			}
		}
	}()

	fmt.Print("Общий чат: \n")

	for {
		reader := bufio.NewReader(os.Stdin)
		source, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Некорректный ввод", err)
			continue
		}
		// отправляем сообщение серверу

		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			fmt.Println(err)
			return
		}

		if strings.TrimSpace(source) == "/exit" {
			break
		}

	}
	wg.Wait()
}
