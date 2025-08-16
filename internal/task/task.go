package task

import "errors"

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Status      Status `json:"status"`
	Description string `json:"description"`
}

// ErrTaskNotFound возвращается, когда задача не найдена
var ErrTaskNotFound = errors.New("task not found")
