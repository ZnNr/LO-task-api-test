package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInMemoryRepository_Create проверяет создание задачи в репозитории
func TestInMemoryRepository_Create(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task := &Task{
		Title:  "Test Task",
		Status: StatusPending,
	}

	err := repo.CreateTask(task)
	require.NoError(t, err, "CreateTask should not return error")

	assert.Equal(t, 1, task.ID, "Expected ID 1")

	retrieved, err := repo.GetTaskByID(1)
	require.NoError(t, err, "GetTaskByID should not return error")

	assert.Equal(t, "Test Task", retrieved.Title, "Expected title 'Test Task'")
	assert.Equal(t, StatusPending, retrieved.Status, "Expected status 'pending'")
}

// TestInMemoryRepository_GetAll проверяет получение всех задач с фильтрацией
func TestInMemoryRepository_GetAll(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	// несколько task
	require.NoError(t, repo.CreateTask(&Task{Title: "Task 1", Status: StatusPending}))
	require.NoError(t, repo.CreateTask(&Task{Title: "Task 2", Status: StatusInProgress}))
	require.NoError(t, repo.CreateTask(&Task{Title: "Task 3", Status: StatusCompleted}))

	// Получим все task
	all, err := repo.GetAllTask(nil)
	require.NoError(t, err, "GetAllTask should not return error")
	assert.Len(t, all, 3, "Expected 3 tasks")

	statusFilter := StatusPending
	pending, err := repo.GetAllTask(&statusFilter)
	require.NoError(t, err, "GetAllTask with filter should not return error")
	assert.Len(t, pending, 1, "Expected 1 pending task")
	assert.Equal(t, "Task 1", pending[0].Title, "Expected 'Task 1'")
}

// TestInMemoryRepository_GetTaskByID проверяет получение задачи по ID
func TestInMemoryRepository_GetTaskByID(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	task := &Task{Title: "Find Me", Status: StatusInProgress}
	require.NoError(t, repo.CreateTask(task))

	found, err := repo.GetTaskByID(task.ID)
	require.NoError(t, err, "GetTaskByID should not return error")
	assert.Equal(t, "Find Me", found.Title, "Expected 'Find Me'")

	_, err = repo.GetTaskByID(999)
	assert.Error(t, err, "Expected error for non-existent ID")
}

// TestInMemoryRepository_ConcurrentAccess проверяет потокобезопасность репозитория
func TestInMemoryRepository_ConcurrentAccess(t *testing.T) {
	repo := NewInMemoryTaskRepository()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			task := &Task{Title: "Concurrent Task", Status: StatusPending}
			_ = repo.CreateTask(task) // Ошибки игнорируем в тестах конкурентности
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	all, _ := repo.GetAllTask(nil)
	assert.Len(t, all, 10, "Expected 10 tasks")
}
