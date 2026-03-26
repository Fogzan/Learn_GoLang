package models

type Task struct {
	ID          int
	Title       string
	Description string
	Completed   bool
}

// Функция для создания новой задачи (конструктор)
func NewTask(id int, title, description string) *Task {
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Completed:   false,
	}
}
