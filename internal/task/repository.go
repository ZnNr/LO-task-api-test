// Package task предоставляет модели и логику работы с задачами
package task

import (
	"sync"
)

// Repository определяет интерфейс для работы с хранилищем задач
type Repository interface {
	CreateTask(task *Task) error
	GetAllTask(status *Status) ([]Task, error)
	GetTaskByID(id int) (*Task, error)
}

type InMemoryTaskRepository struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (r *InMemoryTaskRepository) CreateTask(task *Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.tasks[task.ID] = *task
	r.nextID++
	return nil
}

func (r *InMemoryTaskRepository) GetAllTask(status *Status) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Task
	for _, t := range r.tasks {
		if status == nil || t.Status == *status {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *InMemoryTaskRepository) GetTaskByID(id int) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if t, exists := r.tasks[id]; exists {
		return &t, nil
	}
	return nil, ErrTaskNotFound
}
