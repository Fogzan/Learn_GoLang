package handlers

import (
	_ "fmt"
	"taskmaster/models"
	"taskmaster/storage"
	_ "taskmaster/storage"
)

func CreateTask(title, desc string) *models.Task {
	result := &models.Task{Title: title, Description: desc, Completed: false, ID: storage.GetId()}
	storage.Add(result)
	return result
}

func CompleteTask(id int) {
	result := storage.Get(id)
	result.Completed = true
	storage.Update(result)
}

func ListsTask() []*models.Task {
	return storage.GetAll()
}

func ListPendingTasks() []*models.Task {
	allTasks := storage.GetAll()
	result := []*models.Task{}
	for _, item := range allTasks {
		if item.Completed {
			result = append(result, item)
		}
	}
	return result
}

func ListUnfulfilledTasks() []*models.Task {
	allTasks := storage.GetAll()
	result := []*models.Task{}
	for _, item := range allTasks {
		if !item.Completed {
			result = append(result, item)
		}
	}
	return result
}

func DeleteTask(id int) {
	storage.Delete(id)
}
