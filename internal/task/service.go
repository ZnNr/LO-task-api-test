// Package task предоставляет модели и логику работы с задачами
package task

import (
	"context"
)

// Service определяет интерфейс для работы с задачами
type Service interface {
	CreateTask(ctx context.Context, title string, status Status) error
	GetAllTask(ctx context.Context, status *Status) ([]Task, error)
	GetTaskByID(ctx context.Context, id int) (*Task, error)
}

type TaskService struct {
	repo Repository
}

func NewTaskService(repo Repository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, title string, status Status) error {
	task := &Task{
		Title:  title,
		Status: status,
	}
	return s.repo.CreateTask(task)
}

func (s *TaskService) GetAllTask(ctx context.Context, status *Status) ([]Task, error) {
	return s.repo.GetAllTask(status)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	task, err := s.repo.GetTaskByID(id)
	if err != nil {
		// Проверяем, является ли ошибка "не найдено"
		if err.Error() == "task not found" {
			return nil, ErrTaskNotFound
		}
		// Возвращаем оригинальную ошибку
		return nil, err
	}
	return task, nil
}
