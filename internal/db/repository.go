package db

import (
	"errors"
	"todo/pkg/sqlite3"
)

func TaskRepositoryInit() (*TaskRepository, error) {
	dbConn, err := sqlite3.ConnectorInit()

	if err != nil {
		return nil, err
	}

	dbRepo := &TaskRepository{
		db: dbConn,
	}
	err = dbRepo.createTaskTable()

	if err != nil {
		return nil, err
	}

	return dbRepo, err
}

func (repository *TaskRepository) createTaskTable() error {

	query := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		due_date TEXT,
		completed INTEGER NOT NULL CHECK (completed IN (0, 1)),
		overdue INTEGER NOT NULL CHECK (overdue IN (0, 1)),
		created_at TEXT NOT NULL
	);`

	_, err := repository.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (repository *TaskRepository) CreateTask(input *TaskInput) (*Task, error) {
	result, err := repository.db.Exec("INSERT INTO tasks (title, description, due_date, completed, overdue, created_at) VALUES ($1, $2, $3, 0, 0, $4)",
		input.Title, input.Description, input.DueDate, input.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Получаем ID вставленной задачи
	taskID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	task := &Task{
		ID:          int(taskID),
		Title:       *input.Title,
		Description: *input.Description,
		CreatedAt:   input.CreatedAt,
		DueDate:     *input.DueDate,
		Completed:   0,
		Overdue:     0,
	}

	return task, nil
}

func (repository *TaskRepository) GetAllTasks() ([]*Task, error) {
	rows, err := repository.db.Query("SELECT id, title, description, due_date, completed, overdue, created_at FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Completed, &task.Overdue, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (repository *TaskRepository) GetTaskById(id int) (*Task, error) {
	row := repository.db.QueryRow("SELECT id, title, description, due_date, completed, overdue, created_at FROM tasks WHERE id = $1", id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Completed, &task.Overdue, &task.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (repository *TaskRepository) UpdateTask(task *Task) error {
	result, err := repository.db.Exec("UPDATE tasks SET title = $1, description = $2, due_date = $3, completed = $4, overdue = $5 WHERE id = $6", task.Title, task.Description, task.DueDate, task.Completed, task.Overdue, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (repository *TaskRepository) DeleteTask(taskID int) error {
	result, err := repository.db.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (repository *TaskRepository) CompleteTask(taskID int) error {
	result, err := repository.db.Exec("UPDATE tasks SET completed = 1 WHERE id = $1", taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

// UpdateOverdueTasks обновляет статус просроченных задач
func (repository *TaskRepository) UpdateOverdueTasks(now string) (int64, error) {
	result, err := repository.db.Exec("UPDATE tasks SET overdue = 1 WHERE due_date < $1 AND completed = 0 AND overdue = 0", now)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}
