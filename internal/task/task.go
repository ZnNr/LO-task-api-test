package task

type Status string

const (
	StatusPending    Status = "Pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Status      Status `json:"status"`
	Description string `json:"description"`
}
