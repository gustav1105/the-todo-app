package model

import (
	"time"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"errors"
)

// TaskStatus represents the task status enum in the database
type TaskStatus int32

const (
	TaskStatusTODO      TaskStatus = 0
	TaskStatusDUE       TaskStatus = 1
	TaskStatusOVERDUE   TaskStatus = 2
	TaskStatusCANCELLED TaskStatus = 3
	TaskStatusCOMPLETED TaskStatus = 4
)

// Task represents a task record in the database
type Task struct {
	ID     int64      `db:"id" json:"id"`           // Auto-increment ID
	Name   string     `db:"name" json:"name"`       // Task name
	Note   string     `db:"note" json:"note"`       // Additional notes for the task
	Status TaskStatus `db:"status" json:"status"`   // Status of the task (enum)
	DueOn  *time.Time `db:"due_on" json:"due_on"`   // Due date for the task
}

// TaskModel defines the model that handles operations for tasks
type TaskModel struct {
	DB *sqlx.DB
}

// NewTaskModel creates a new instance of TaskModel
func NewTaskModel(db *sqlx.DB) *TaskModel {
	return &TaskModel{DB: db}
}

// CreateTask inserts a new task into the database
func (m *TaskModel) CreateTask(task *Task) error {
	query := `INSERT INTO tasks (name, note, status, due_on) VALUES (?, ?, ?, ?)`
	_, err := m.DB.Exec(query, task.Name, task.Note, task.Status, task.DueOn)
	if err != nil {
		return err
	}
	return nil
}

// GetTask retrieves a task by its ID
func (m *TaskModel) GetTask(id int64) (*Task, error) {
	var task Task
	query := `SELECT id, name, note, status, due_on FROM tasks WHERE id = ?`
	err := m.DB.Get(&task, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Task not found
		}
		return nil, err
	}
	return &task, nil
}

// GetTasks retrieves all tasks from the database
func (m *TaskModel) GetTasks() ([]Task, error) {
	var tasks []Task
	query := `SELECT id, name, note, status, due_on FROM tasks`
	err := m.DB.Select(&tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// UpdateTaskStatus updates the status of a task by ID
func (m *TaskModel) UpdateTaskStatus(id int64, status TaskStatus) error {
	query := `UPDATE tasks SET status = ? WHERE id = ?`
	_, err := m.DB.Exec(query, status, id)
	return err
}

// DeleteTask removes a task by its ID
func (m *TaskModel) DeleteTask(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := m.DB.Exec(query, id)
	return err
}

