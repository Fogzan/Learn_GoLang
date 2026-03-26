package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const filenameConstant string = "contacts.txt"

type Contact struct {
	Name  string
	Phone string
	Email string
}

func SaveContact(filename string, contacts []Contact) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, item := range contacts {
		fmt.Fprintf(file, "%s|%s|%s\n", item.Name, item.Phone, item.Email)
	}
	return nil
}

func LoadContact(filename string) ([]Contact, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []Contact
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue // или вернуть ошибку
		}
		result = append(result, Contact{
			Name:  parts[0],
			Phone: parts[1],
			Email: parts[2],
		})
	}
	return result, scanner.Err()
}

func AddContact(contacts *[]Contact) {
	var name, phone, email string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Введите имя: ")
	scanner.Scan()
	name = scanner.Text()
	flag := true
	for {
		flag = true
		fmt.Printf("Введите телефон: ")
		scanner.Scan()
		phone = scanner.Text()
		for _, elem := range phone {
			if !(elem >= '0' && elem <= '9') {
				flag = false
				fmt.Printf("Ошибка! Это не телефон!\n")
				break
			}
		}
		if flag {
			break
		}
	}

	for {
		fmt.Printf("Введите почту: ")
		scanner.Scan()
		email = scanner.Text()
		if !strings.Contains(email, "@") {
			fmt.Printf("Ошибка! Это не почта!\n")
			continue
		}
		break
	}

	*contacts = append(*contacts, Contact{Name: name, Phone: phone, Email: email})
}

func SearchContact(contacts []Contact, name string) (*Contact, error) {
	for i, item := range contacts {
		if item.Name == name {
			return &contacts[i], nil
		}
	}
	return nil, errors.New("Такой контакт не найден")
}

func mainMenu(choise int, contacts *[]Contact) error {
	switch choise {
	case 1:
		for i, item := range *contacts {
			fmt.Printf("Контакт №%v: имя-%v телефон-%v почта-%v\n", i+1, item.Name, item.Phone, item.Email)
		}
		return nil
	case 2:
		AddContact(contacts)
		return nil
	case 3:
		fmt.Printf("Введите имя для поиска: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		name := scanner.Text()

		finditem, err := SearchContact(*contacts, name)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		} else {
			fmt.Printf("Найден контакт: имя-%v телефон-%v почта-%v\n", finditem.Name, finditem.Phone, finditem.Email)
		}
		return nil
	case 4:
		err := SaveContact(filenameConstant, *contacts)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		}
		return nil
	case 5:
		result, err := LoadContact(filenameConstant)
		if err != nil {
			fmt.Printf("Ошибка при первом чтении файла: %v\n", err)
		} else {
			*contacts = result
		}
		return nil
	default:
		return errors.New("Ошибочный выбор")
	}
}

func main() {
	contacts, err := LoadContact(filenameConstant)
	if err != nil {
		fmt.Printf("Ошибка при первом чтении файла: %v\n", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	var choise int
	for {
		fmt.Printf("\n\nВыберите действие: \n1. Показать все контакты\n2. Добавить контакт\n3. Найти контакт по имени\n4. Сохранить в файл\n5. Загрузить из файла\n6. Выход\n\t>>> ")
		scanner.Scan()
		input := scanner.Text()
		choise, err = strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		}
		if choise == 6 {
			break
		} else {
			err := mainMenu(choise, &contacts)
			if err != nil {
				fmt.Printf("Ошибка при выборе меню: %v\n", err)
			}
		}
	}

}
