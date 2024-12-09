package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo/internal/db"
)

func TestGetTasks(t *testing.T) {
	mockRepo := NewMockRepository()
	handler := &Handler{repo: mockRepo}

	// Добавляем тестовые задачи в мок-репозиторий
	mockRepo.tasks[0] = db.Task{ID: 1, Title: "Task 1", DueDate: "2006-01-02 15:55:08", Completed: 0, Overdue: 0, CreatedAt: "2006-01-02 15:45:08"}
	mockRepo.tasks[1] = db.Task{ID: 2, Title: "Task 2", DueDate: "2006-01-02 15:55:08", Completed: 0, Overdue: 0, CreatedAt: "2006-01-02 15:45:08"}

	// Эмулируем запрос
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.getTasks(rr, req)

	// Проверяем статус
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что тело ответа соответствует ожидаемому

	var gotTasks []db.Task
	if err := json.NewDecoder(rr.Body).Decode(&gotTasks); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Ожидаемый результат
	wantTasks := []db.Task{
		{ID: 1, Title: "Task 1", DueDate: "2006-01-02 15:55:08", Completed: 0, Overdue: 0, CreatedAt: "2006-01-02 15:45:08"},
		{ID: 2, Title: "Task 2", DueDate: "2006-01-02 15:55:08", Completed: 0, Overdue: 0, CreatedAt: "2006-01-02 15:45:08"},
	}

	// Сравниваем структуры
	if !equalTasks(gotTasks, wantTasks) {
		t.Errorf("handler returned unexpected body: got %+v, want %+v", gotTasks, wantTasks)
	}
}

func equalTasks(got, want []db.Task) bool {
	if len(got) != len(want) {
		return false
	}

	for i := range got {
		if got[i].ID != want[i].ID || got[i].Title != want[i].Title ||
			got[i].DueDate != want[i].DueDate || got[i].Completed != want[i].Completed ||
			got[i].Overdue != want[i].Overdue || got[i].CreatedAt != want[i].CreatedAt {
			return false
		}
	}
	return true
}

func TestCreateTask(t *testing.T) {
	mockRepo := &MockRepository{}
	handler := &Handler{repo: mockRepo}

	title := "New Task"

	taskInput := db.TaskInput{
		Title: &title,
	}

	// Создаем тело запроса
	body, err := json.Marshal(taskInput)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.createTask(rr, req)

	var gotTask db.Task
	if err := json.NewDecoder(rr.Body).Decode(&gotTask); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Проверяем статус
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var gotArr, mockArr []db.Task
	gotArr = append(gotArr, gotTask)
	mockArr = append(mockArr, mockRepo.tasks[0])
	if !equalTasks(gotArr, mockArr) {
		t.Errorf("handler returned unexpected body: got %v want %v", gotTask, mockRepo.tasks[0])
	}
}

func TestUpdateTask(t *testing.T) {
	mockRepo := NewMockRepository()
	handler := &Handler{repo: mockRepo}

	mockRepo.tasks[0] = db.Task{ID: 1, Title: "Task 1", DueDate: "2006-01-02 15:55:08", Completed: 0, Overdue: 0, CreatedAt: "2006-01-02 15:45:08"}

	title := "Updated Task"

	updatedTaskInput := db.TaskInput{
		Title: &title,
	}

	// Создаем тело запроса
	body, err := json.Marshal(updatedTaskInput)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.updateTask(rr, req, 1)

	// Проверяем статус
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if *updatedTaskInput.Title != mockRepo.tasks[0].Title {
		t.Errorf("handler returned unexpected body: got %v want %v", mockRepo.tasks[0].Title, *updatedTaskInput.Title)
	}
}

func TestDeleteTask(t *testing.T) {
	mockRepo := NewMockRepository()
	handler := &Handler{repo: mockRepo}

	req, err := http.NewRequest("DELETE", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.deleteTask(rr, req, 1)

	// Проверяем статус
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	mockRepo.tasks[0] = db.Task{ID: 1, Title: "Task 1", DueDate: "2006-01-02 15:55:08", Completed: 0, Overdue: 0, CreatedAt: "2006-01-02 15:45:08"}

	req, err = http.NewRequest("DELETE", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.deleteTask(rr, req, 0)

	// Проверяем статус
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestCompleteTask(t *testing.T) {
	mockRepo := NewMockRepository()
	handler := &Handler{repo: mockRepo}

	mockRepo.tasks[1] = db.Task{ID: 1, Title: "Task 1", DueDate: "2006-01-02 15:55:08", Completed: 0, Overdue: 0, CreatedAt: "2006-01-02 15:45:08"}

	req, err := http.NewRequest("PATCH", "/tasks/1/complete", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.completeTask(rr, req, 1)

	// Проверяем статус
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if mockRepo.tasks[1].Completed != 1 {
		t.Errorf("handler returned unexpected body: got %v want %v", mockRepo.tasks[0].Completed, 1)
	}
}
