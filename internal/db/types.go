package db

import (
	"database/sql"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Completed   int8   `json:"completed"`
	Overdue     int8   `json:"overdue"`
	CreatedAt   string `json:"created_at"`
}

type DbInterface interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

type TaskRepository struct {
	db DbInterface
}

type TaskInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
}
