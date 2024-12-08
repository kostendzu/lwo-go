package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"todo/internal/db"

	"math/rand"
)

func ifEmptyUseCurrent(updatedValue *string, currentValue string) string {
	if updatedValue == nil {
		return currentValue
	}
	return *updatedValue
}

func validateTaskInput(input *db.TaskInput) error {
	if input.Title == nil {
		return fmt.Errorf("title is required")
	}

	// Проверка на правильность формата даты (если указана)
	if input.DueDate != nil {
		if _, err := time.Parse("2006-01-02", *input.DueDate); err != nil {
			return fmt.Errorf("invalid due_date format, expected YYYY-MM-DD")
		}
	}

	return nil
}

func transformTaskInput(input *db.TaskInput) *db.TaskInput {
	if input.DueDate != nil {
		*input.DueDate += " 00:00:00"
	} else {
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		randomNumber := r.Intn(11)
		dueDate := time.Now().Add(time.Duration(randomNumber) * time.Minute).Format("2000-01-01 00:00:00")
		input.DueDate = &dueDate
	}

	return input
}

func checkDueDate(dueDate string, createdAt string) error {
	dueDateTime, _ := time.Parse("2006-01-02 00:00:00", dueDate)
	createdAtTime, _ := time.Parse("2006-01-02 15:04:05", createdAt)

	if createdAtTime.Compare(dueDateTime) >= 0 {
		return errors.New("DueDate <= createdAt")
	}

	return nil
}

// GET /tasks - Получить все задачи
func (h *Handler) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.GetAllTasks()
	if err != nil {
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// POST /tasks - Создать новую задачу
func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var input *db.TaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	input.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	if err := validateTaskInput(input); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	input = transformTaskInput(input)

	if err := checkDueDate(*input.DueDate, input.CreatedAt); err != nil {
		http.Error(w, fmt.Sprintf("Invalid dueDate: %v", err), http.StatusBadRequest)
		return
	}

	task, err := h.repo.CreateTask(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create task: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// PUT /tasks/{id} - Обновить задачу
func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request, id int) {
	currentTask, err := h.repo.GetTaskById(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid task id: %v", err), http.StatusBadRequest)
		return
	}

	var updatedTaskInput *db.TaskInput
	if err := json.NewDecoder(r.Body).Decode(&updatedTaskInput); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if err := checkDueDate(*updatedTaskInput.DueDate, currentTask.CreatedAt); err != nil {
		http.Error(w, fmt.Sprintf("Invalid dueDate: %v", err), http.StatusBadRequest)
		return
	}

	var updatedTask = &db.Task{
		ID:          currentTask.ID,
		Title:       ifEmptyUseCurrent(updatedTaskInput.Title, currentTask.Title),
		Description: ifEmptyUseCurrent(updatedTaskInput.Description, currentTask.Description),
		DueDate:     ifEmptyUseCurrent(updatedTaskInput.DueDate, currentTask.DueDate),
		Completed:   currentTask.Completed,
		Overdue:     currentTask.Overdue,
		CreatedAt:   currentTask.CreatedAt,
	}

	err = h.repo.UpdateTask(updatedTask)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}

// DELETE /tasks/{id} - Удалить задачу
func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	err := h.repo.DeleteTask(id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PATCH /tasks/{id}/complete - Завершить задачу
func (h *Handler) completeTask(w http.ResponseWriter, r *http.Request, id int) {
	err := h.repo.CompleteTask(id)
	if err != nil {
		http.Error(w, "Failed to mark task as completed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateOverdueTasks() (int64, error) {
	now := time.Now().Format("2006-01-02")
	updatedCount, err := h.repo.UpdateOverdueTasks(now)
	if err != nil {
		return 0, err
	}

	return updatedCount, nil
}
