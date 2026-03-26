package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"taskmaster/handlers"
)

func demo(status string) {
	fmt.Println("========================================================")
	list := handlers.ListsTask()
	switch status {
	case "true":
		list = handlers.ListPendingTasks()
	case "false":
		list = handlers.ListUnfulfilledTasks()
	}
	for _, item := range list {
		fmt.Println(*item)
	}
	fmt.Println("========================================================")
}

func menu() {
	scanner := bufio.NewScanner(os.Stdin)
mainfor:
	for true {
		fmt.Println("\nВыберите действие: ")
		fmt.Println("1. Добавить задачу")
		fmt.Println("2. Отметить задачу выполненной")
		fmt.Println("3. Показать все задачи")
		fmt.Println("4. Показать невыполненные задачи")
		fmt.Println("5. Удалить задачу")
		fmt.Println("6. Выход")
		fmt.Printf("> ")
		scanner.Scan()
		idMenu := scanner.Text()
		switch idMenu {
		case "1":
			fmt.Printf("Введите название: ")
			scanner.Scan()
			title := scanner.Text()
			fmt.Printf("Введите описание: ")
			scanner.Scan()
			desc := scanner.Text()
			handlers.CreateTask(title, desc)
		case "2":
			fmt.Printf("Все задачи: \n")
			demo("all")
			fmt.Printf("Введите id задачи >")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())
			handlers.CompleteTask(id)
		case "3":
			fmt.Printf("Все задачи: \n")
			demo("all")
		case "4":
			fmt.Printf("Невыполненные задачи: \n")
			demo("false")
		case "5":
			fmt.Printf("Все задачи: \n")
			demo("all")
			fmt.Printf("Введите id задачи >")
			scanner.Scan()
			id, _ := strconv.Atoi(scanner.Text())
			handlers.DeleteTask(id)
		case "6":
			break mainfor
		default:
			fmt.Println("Неизвестная команда")
		}
	}
}

func main() {
	menu()
}
