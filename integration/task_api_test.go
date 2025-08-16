// Package integration предоставляет интеграционные тесты для API задач
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ZnNr/LO-task-api-test/internal/handler"
	"github.com/ZnNr/LO-task-api-test/internal/task"
)

// TestTaskAPI_Integration тестирует основные сценарии работы с задачами
func TestTaskAPI_Integration(t *testing.T) {

	repo := task.NewInMemoryTaskRepository()
	service := task.NewTaskService(repo)
	logChan := make(chan string, 100)

	// Асинхронное логирование
	go func() {
		for range logChan {
			// игнор сообщений в тестах, но канал читаем
			// чтобы не блокировать логгер
		}
	}()

	taskHandler := handler.NewTaskHandler(service, logChan)

	// Создание задачи
	t.Run("CreateTask", func(t *testing.T) {
		body := `{"title":"Integration Test Task","status":"pending"}`
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		taskHandler.HandleTasks(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Expected 201 Created")
	})

	// Получить все уже созданные задачи
	t.Run("GetAllTasks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()

		taskHandler.HandleTasks(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

		var tasks []task.Task
		err := json.Unmarshal(w.Body.Bytes(), &tasks)
		require.NoError(t, err, "Should unmarshal JSON")

		assert.Len(t, tasks, 1, "Expected 1 task")
		assert.Equal(t, "Integration Test Task", tasks[0].Title, "Expected correct title")
	})

	// Получение задачи по ID
	t.Run("GetTaskByID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
		w := httptest.NewRecorder()

		taskHandler.HandleTaskByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

		var taskResp task.Task
		err := json.Unmarshal(w.Body.Bytes(), &taskResp)
		require.NoError(t, err, "Should unmarshal JSON")

		assert.Equal(t, "Integration Test Task", taskResp.Title, "Expected correct title")
	})

	// Фильтрация по статусу
	t.Run("GetTasksWithFilter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks?status=pending", nil)
		w := httptest.NewRecorder()

		taskHandler.HandleTasks(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

		var tasks []task.Task
		err := json.Unmarshal(w.Body.Bytes(), &tasks)
		require.NoError(t, err, "Should unmarshal JSON")

		assert.Len(t, tasks, 1, "Expected 1 pending task")
	})

	// Задача не найдена
	t.Run("GetTaskByID_NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
		w := httptest.NewRecorder()

		taskHandler.HandleTaskByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Expected 404 Not Found")
	})

	// Проверяет на валидный JSON
	t.Run("CreateTask_InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{invalid json`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		taskHandler.HandleTasks(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Expected 400 Bad Request")
	})
}

// TestTaskAPI_Concurrent проверяет, что API корректно обрабатывает
// одновременное создание нескольких задач без возникновения гонок или ошибок
func TestTaskAPI_Concurrent(t *testing.T) {
	repo := task.NewInMemoryTaskRepository()
	service := task.NewTaskService(repo)
	logChan := make(chan string, 1000)

	go func() {
		for range logChan {
		}
	}()

	taskHandler := handler.NewTaskHandler(service, logChan)

	// Создаем 10 задач параллельно
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			body := `{"title":"Concurrent Task ` + string(rune('0'+i)) + `","status":"pending"}`
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			taskHandler.HandleTasks(w, req)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent requests")
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	taskHandler.HandleTasks(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

	var tasks []task.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	require.NoError(t, err, "Should unmarshal JSON")

	assert.Len(t, tasks, 10, "Expected 10 tasks")
}

// TestTaskAPI_ContextTimeout проверяет, что API корректно
// обрабатывает HTTP-запросы с ограниченным по времени контекстом
func TestTaskAPI_ContextTimeout(t *testing.T) {
	repo := task.NewInMemoryTaskRepository()
	service := task.NewTaskService(repo)
	logChan := make(chan string, 100)

	go func() {
		for range logChan {
		}
	}()

	taskHandler := handler.NewTaskHandler(service, logChan)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	taskHandler.HandleTasks(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")
}
