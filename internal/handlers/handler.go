package handlers

import (
	"net/http"
	"strconv"
	"todo/internal/db"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	repo db.Repo
}

func NewHandler(repo *db.TaskRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getTasks(w, r)
	case "POST":
		h.createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/tasks/"):])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		h.updateTask(w, r, id)
	case "DELETE":
		h.deleteTask(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/tasks/complete/"):])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if r.Method == "PATCH" {
		h.completeTask(w, r, id)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
