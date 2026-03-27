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
	// defer conn.Close() перенес в горутину

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				wg.Done()
				break
			}
			if string(buff[0:n]) != "true" {
				fmt.Print(string(buff[0:n]))
			} else {
				conn.Close()
				wg.Done()
			}
		}
	}()

	reader := bufio.NewReader(os.Stdin)
InputUsernameFor:
	for {
		fmt.Printf("Введите никнейм: \n> ")
		myUserName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Некорректный ввод", err)
		}
		if strings.TrimSpace(myUserName) == "" {
			fmt.Printf("Некорректный ввод\n")
		} else {
			n, err := conn.Write([]byte(myUserName))
			if n == 0 || err != nil {
				fmt.Println(err)
				return
			}
			break InputUsernameFor
		}
	}

	fmt.Printf("\nОбщий чат: \n")

	for {
		source, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Некорректный ввод", err)
			continue
		}
		// отправляем сообщение серверу
		if n, err := conn.Write([]byte(source)); n == 0 || err != nil {
			fmt.Println(err)
			continue
		}

		if strings.TrimSpace(source) == "/exit" {
			// buff := make([]byte, 1024)
			// n, err := conn.Read(buff)
			// if err != nil {
			// 	break
			// }
			// if string(buff[0:n]) == "true" {
			// 	break
			// } else {
			// 	continue
			// }
			break
		}

	}
	wg.Wait()
}
