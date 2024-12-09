package handlers

import (
	"fmt"
	"time"
	"todo/internal/db"
)

type MockRepository struct {
	tasks map[int]db.Task
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		tasks: make(map[int]db.Task),
	}
}

func (m *MockRepository) GetAllTasks() ([]*db.Task, error) {
	var result []*db.Task
	for _, task := range m.tasks {
		result = append(result, &task)
	}
	return result, nil
}
func (m *MockRepository) CreateTask(input *db.TaskInput) (*db.Task, error) {
	// Инициализируем карту, если она еще не была инициализирована
	if m.tasks == nil {
		m.tasks = make(map[int]db.Task)
	}

	// Пример создания задачи (вы можете изменить логику, если нужно)
	task := db.Task{
		ID:        len(m.tasks),
		Title:     *input.Title,
		DueDate:   *input.DueDate,
		Completed: 0,
		Overdue:   0,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Добавляем задачу в карту
	m.tasks[task.ID] = task

	return &task, nil
}

func (m *MockRepository) GetTaskById(id int) (*db.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}
	return &task, nil
}

func (m *MockRepository) UpdateTask(task *db.Task) error {
	_, exists := m.tasks[task.ID]
	if !exists {
		return fmt.Errorf("task not found")
	}
	m.tasks[task.ID] = *task
	return nil
}

func (m *MockRepository) DeleteTask(id int) (int64, error) {
	_, exists := m.tasks[id]
	if !exists {
		return 0, nil
	}
	delete(m.tasks, id)
	return 1, nil
}

func (m *MockRepository) CompleteTask(id int) error {
	task, exists := m.tasks[id]
	if !exists {
		return fmt.Errorf("task not found")
	}
	task.Completed = 1
	m.tasks[task.ID] = task
	return nil
}

func (m *MockRepository) UpdateOverdueTasks(currentTime string) (int64, error) {
	// Возвращаем заранее заданные данные
	return 0, nil
}
