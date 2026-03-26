package storage

import (
	"fmt"
	"taskmaster/models"
)

var storage struct {
	tasks  []*models.Task
	nextId int
}

func Add(task *models.Task) {
	storage.tasks = append(storage.tasks, task)
	storage.nextId++
}

func GetId() int {
	return storage.nextId
}

func findTask(id int) (*models.Task, int) {
	for i, item := range storage.tasks {
		if item.ID == id {
			return item, i
		}
	}
	return nil, -1
}

func Get(id int) *models.Task {
	item, _ := findTask(id)
	if item == nil {
		fmt.Println("Ошибка пакета storage: Не существующий id")
		return nil
	}
	return item
}

func GetAll() []*models.Task {
	return storage.tasks
}

func Update(task *models.Task) {
	oldTask, i := findTask(task.ID)
	if oldTask == nil {
		fmt.Println("Ошибка пакета storage: Не существующий id")
		return
	}
	storage.tasks[i] = task
}

func Delete(id int) {
	oldTask, i := findTask(id)
	if oldTask == nil {
		fmt.Println("Ошибка пакета storage: Не существующий id")
		return
	}
	storage.tasks = append(storage.tasks[:i], storage.tasks[i+1:]...)
}
