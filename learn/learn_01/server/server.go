package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type connectedUser struct {
	socket   net.Conn
	address  net.Addr
	id       int
	userName string
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

	var mute sync.Mutex
	connectedClient := []connectedUser{}
	idClient := 0

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			conn.Close()
			return
		}
		go processingConnection(conn, &connectedClient, idClient, &mute)
		idClient++
	}
}

func processingConnection(conn net.Conn, connectedClient *[]connectedUser, idClient int, mute *sync.Mutex) {
	// Доп часть
	defer conn.Close()
	defer func() {
		mute.Lock()
		fmt.Printf("Количество клиентов: %v\n\n", len(*connectedClient))
		mute.Unlock()
	}()
	defer fmt.Printf("\n[Пользователь %v] Отключился\n", conn.RemoteAddr())
	fmt.Printf("\n[Пользователь %v] Подключился\n", conn.RemoteAddr())

	// Получение username
	input := make([]byte, (1024 * 4))
	n, err := conn.Read(input)
	if n == 0 || err != nil {
		fmt.Println("Error: ", err)
		return
	}
	getUserName := strings.TrimSpace(string(input[:n]))
	fmt.Printf("Пользователь: %v - назвался: %v", conn.RemoteAddr(), getUserName)

	// Добавление клиента в массив
	mute.Lock()
	*connectedClient = append(*connectedClient, connectedUser{socket: conn, address: conn.RemoteAddr(), id: idClient, userName: getUserName})
	mute.Unlock()

	// функция которая выполняется после разрыва соеденения с клиентом
	defer func() {
		mute.Lock()
		defer mute.Unlock()
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

	// Вывод количества клиентов
	mute.Lock()
	fmt.Printf("Количество клиентов: %v\n\n", len(*connectedClient))
	mute.Unlock()

	// Основной цикл
	for {
		// Получение сообщения
		input := make([]byte, (1024 * 4))
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			fmt.Println("Error: ", err)
			return
		}
		sms := string(input[:n])
		if strings.TrimSpace(sms) != "" {
			fmt.Printf("Пользователь: %v - оставил сообщение: %v", conn.RemoteAddr(), sms)
		}

		// Обработка полученного сообщения
		switch strings.TrimSpace(sms) {
		case "/exit":
			_, err := conn.Write([]byte("true"))
			if err != nil {
				fmt.Println(err)
			}
			return
		case "/list":
			clientList := "Пользователи в чате: "
			if len(*connectedClient) < 2 {
				clientList = "Ты в чате один!\n"
			} else {
				mute.Lock()
				for _, item := range *connectedClient {
					if item.id != idClient {
						clientList += item.userName + " "
					}
				}
				mute.Unlock()
				clientList += "\n"
			}

			_, err := conn.Write([]byte(clientList))
			if err != nil {
				fmt.Println(err)
			}
		case "":
			continue
		default:
			mute.Lock()
			if len(*connectedClient) < 2 {
				mute.Unlock()
				continue
			}
			for _, item := range *connectedClient {
				if item.id != idClient {
					item.socket.SetWriteDeadline(time.Now().Add(30 * time.Second))
					_, err := item.socket.Write([]byte("\t\t" + getUserName + ": " + sms))
					if err != nil {
						fmt.Println(err)
					}
					item.socket.SetWriteDeadline(time.Time{})
				}
			}
			mute.Unlock()
		}

	}

}
