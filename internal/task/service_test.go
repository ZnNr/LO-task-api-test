// Package task предоставляет модели и логику работы с задачами
package task

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTaskService_Create проверяет успешное создание задачи через сервис
func TestTaskService_Create(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo)

	err := service.CreateTask(context.Background(), "Test Task", StatusPending)
	require.NoError(t, err, "CreateTask should not return error")

	tasks, err := mockRepo.GetAllTask(nil)
	require.NoError(t, err, "GetAllTask should not return error")
	require.Len(t, tasks, 1, "Expected 1 task")
	assert.Equal(t, "Test Task", tasks[0].Title, "Expected title 'Test Task'")
}

// TestTaskService_Create_Error проверяет обработку ошибок при создании задачи
func TestTaskService_Create_Error(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo)

	mockRepo.SetError(errors.New("database error"))

	err := service.CreateTask(context.Background(), "Test Task", StatusPending)
	assert.Error(t, err, "Expected error")
	assert.Equal(t, "database error", err.Error(), "Expected 'database error'")
}

// TestTaskService_GetAll проверяет получение всех задач через сервис
func TestTaskService_GetAll(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo)

	require.NoError(t, mockRepo.CreateTask(&Task{Title: "Task 1", Status: StatusPending}))
	require.NoError(t, mockRepo.CreateTask(&Task{Title: "Task 2", Status: StatusInProgress}))

	tasks, err := service.GetAllTask(context.Background(), nil)
	require.NoError(t, err, "GetAllTask should not return error")
	assert.Len(t, tasks, 2, "Expected 2 tasks")
}

// TestTaskService_GetAll_WithFilter проверяет фильтрацию задач по статусу
func TestTaskService_GetAll_WithFilter(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo)

	require.NoError(t, mockRepo.CreateTask(&Task{Title: "Task 1", Status: StatusPending}))
	require.NoError(t, mockRepo.CreateTask(&Task{Title: "Task 2", Status: StatusInProgress}))

	statusFilter := StatusPending
	tasks, err := service.GetAllTask(context.Background(), &statusFilter)
	require.NoError(t, err, "GetAllTask with filter should not return error")
	assert.Len(t, tasks, 1, "Expected 1 task")
	assert.Equal(t, "Task 1", tasks[0].Title, "Expected 'Task 1'")
}

// TestTaskService_GetAll_Error проверяет обработку ошибок при получении задач
func TestTaskService_GetAll_Error(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo)

	mockRepo.SetError(errors.New("database error"))

	_, err := service.GetAllTask(context.Background(), nil)
	assert.Error(t, err, "Expected error")
	assert.Equal(t, "database error", err.Error(), "Expected 'database error'")
}

// TestTaskService_GetTaskByID проверяет получение задачи по ID через сервис
func TestTaskService_GetTaskByID(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo)

	require.NoError(t, mockRepo.CreateTask(&Task{Title: "Find Me", Status: StatusCompleted}))

	task, err := service.GetTaskByID(context.Background(), 1)
	require.NoError(t, err, "GetTaskByID should not return error")
	assert.Equal(t, "Find Me", task.Title, "Expected 'Find Me'")
}

// TestTaskService_GetTaskByID_NotFound проверяет обработку случая, когда задача не найдена
func TestTaskService_GetTaskByID_NotFound(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo)

	_, err := service.GetTaskByID(context.Background(), 999)
	assert.Error(t, err, "Expected error for not found")
	assert.True(t, errors.Is(err, ErrTaskNotFound), "Expected ErrTaskNotFound")
}

// TestTaskService_GetTaskByID_Error проверяет обработку ошибок при получении задачи по ID
func TestTaskService_GetTaskByID_Error(t *testing.T) {
	mockRepo := NewMockRepository()
	service := NewTaskService(mockRepo) // ← Исправлено имя функции

	// Установим ошибку
	expectedErr := errors.New("database error")
	mockRepo.SetError(expectedErr)

	_, err := service.GetTaskByID(context.Background(), 1) // ← Исправлено имя метода
	assert.Error(t, err, "Expected error")
	assert.Equal(t, expectedErr, err, "Expected database error")
}
