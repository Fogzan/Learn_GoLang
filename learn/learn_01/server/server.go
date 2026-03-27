package main

import (
	"fmt"
	"net"
	"strings"
)

type connectedUser struct {
	socket  net.Conn
	address net.Addr
	id      int
}

func main() {
	fmt.Println("Server start")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")

	connectedClient := []connectedUser{}
	idClient := 0

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			conn.Close()
			return
		}
		go processingConnection(conn, &connectedClient, idClient)
		idClient++
	}
}

func processingConnection(conn net.Conn, connectedClient *[]connectedUser, idClient int) {
	defer conn.Close()
	defer fmt.Printf("Количество клиентов: %v\n\n", len(*connectedClient))
	defer fmt.Printf("\n[Пользователь %v] Отключился\n", conn.LocalAddr())
	fmt.Printf("\n[Пользователь %v] Подключился\n", conn.LocalAddr())

	*connectedClient = append(*connectedClient, connectedUser{socket: conn, address: conn.LocalAddr(), id: idClient})
	defer func() {
		res := -1
		for n, item := range *connectedClient {
			if item.id == idClient {
				res = n
			}
		}
		if res == -1 {
			return
		} else {
			*connectedClient = append((*connectedClient)[:res], (*connectedClient)[res+1:]...)
		}
	}()

	fmt.Printf("Количество клиентов: %v\n\n", len(*connectedClient))

	for {
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Error: ", err)
			return
		}
		sms := string(input[:n])
		fmt.Printf("Пользователь: %v - оставил сообщение: %v", conn.LocalAddr(), sms)

		if strings.TrimSpace(sms) == "/exit" {
			_, err := conn.Write([]byte("true"))
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		for _, item := range *connectedClient {
			if item.id != idClient {
				_, err := item.socket.Write([]byte(sms))
				if err != nil {
					fmt.Println(err)
				}
			}

		}
	}

}
