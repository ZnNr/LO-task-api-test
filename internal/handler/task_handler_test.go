package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZnNr/LO-task-api-test/internal/task"
)

// MockService реализует интерфейс task.Service для тестирования
type MockService struct {
	GetAllTaskFunc  func(ctx context.Context, status *task.Status) ([]task.Task, error)
	GetTaskByIDFunc func(ctx context.Context, id int) (*task.Task, error)
	CreateTaskFunc  func(ctx context.Context, title string, status task.Status) error
}

func (m *MockService) GetAllTask(ctx context.Context, status *task.Status) ([]task.Task, error) {
	if m.GetAllTaskFunc != nil {
		return m.GetAllTaskFunc(ctx, status)
	}
	return nil, nil
}

func (m *MockService) GetTaskByID(ctx context.Context, id int) (*task.Task, error) {
	if m.GetTaskByIDFunc != nil {
		return m.GetTaskByIDFunc(ctx, id)
	}
	return nil, task.ErrTaskNotFound
}

func (m *MockService) CreateTask(ctx context.Context, title string, status task.Status) error {
	if m.CreateTaskFunc != nil {
		return m.CreateTaskFunc(ctx, title, status)
	}
	return nil
}

func TestTaskHandler_GetTasks_Success(t *testing.T) {
	mockService := &MockService{
		GetAllTaskFunc: func(ctx context.Context, status *task.Status) ([]task.Task, error) {
			return []task.Task{
				{ID: 1, Title: "Task 1", Status: task.StatusPending},
				{ID: 2, Title: "Task 2", Status: task.StatusInProgress},
			}, nil
		},
	}

	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tasks []task.Task
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Проверим логирование
	select {
	case logMsg := <-logChan:
		if logMsg == "" {
			t.Error("Expected log message")
		}
	default:
		t.Error("No log message received")
	}
}

func TestTaskHandler_GetTasks_WithFilter(t *testing.T) {
	mockService := &MockService{
		GetAllTaskFunc: func(ctx context.Context, status *task.Status) ([]task.Task, error) {
			if status != nil && *status == task.StatusPending {
				return []task.Task{
					{ID: 1, Title: "Task 1", Status: task.StatusPending},
				}, nil
			}
			return []task.Task{}, nil
		},
	}

	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodGet, "/tasks?status=pending", nil)
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tasks []task.Task
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestTaskHandler_GetTasks_ServiceError(t *testing.T) {
	mockService := &MockService{
		GetAllTaskFunc: func(ctx context.Context, status *task.Status) ([]task.Task, error) {
			return nil, errors.New("service error")
		},
	}

	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestTaskHandler_CreateTask_Success(t *testing.T) {
	mockService := &MockService{
		CreateTaskFunc: func(ctx context.Context, title string, status task.Status) error {
			return nil
		},
	}

	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	body := `{"title":"New Task","status":"pending"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Проверим логирование
	select {
	case logMsg := <-logChan:
		if logMsg == "" {
			t.Error("Expected log message")
		}
	default:
		t.Error("No log message received")
	}
}

func TestTaskHandler_CreateTask_InvalidJSON(t *testing.T) {
	mockService := &MockService{}
	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestTaskHandler_CreateTask_MissingTitle(t *testing.T) {
	mockService := &MockService{}
	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	body := `{"status":"pending"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestTaskHandler_CreateTask_ServiceError(t *testing.T) {
	mockService := &MockService{
		CreateTaskFunc: func(ctx context.Context, title string, status task.Status) error {
			return errors.New("service error")
		},
	}

	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	body := `{"title":"New Task","status":"pending"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestTaskHandler_GetTaskByID_Success(t *testing.T) {
	mockService := &MockService{
		GetTaskByIDFunc: func(ctx context.Context, id int) (*task.Task, error) {
			return &task.Task{ID: 1, Title: "Test Task", Status: task.StatusPending}, nil
		},
	}

	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	w := httptest.NewRecorder()

	handler.HandleTaskByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var taskResp task.Task
	if err := json.Unmarshal(w.Body.Bytes(), &taskResp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if taskResp.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", taskResp.Title)
	}

	// Проверим логирование
	select {
	case logMsg := <-logChan:
		if logMsg == "" {
			t.Error("Expected log message")
		}
	default:
		t.Error("No log message received")
	}
}

func TestTaskHandler_GetTaskByID_InvalidID(t *testing.T) {
	mockService := &MockService{}
	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodGet, "/tasks/invalid", nil)
	w := httptest.NewRecorder()

	handler.HandleTaskByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestTaskHandler_GetTaskByID_NotFound(t *testing.T) {
	mockService := &MockService{
		GetTaskByIDFunc: func(ctx context.Context, id int) (*task.Task, error) {
			return nil, task.ErrTaskNotFound
		},
	}

	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodGet, "/tasks/999", nil)
	w := httptest.NewRecorder()

	handler.HandleTaskByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestTaskHandler_InvalidMethod(t *testing.T) {
	mockService := &MockService{}
	logChan := make(chan string, 10)
	handler := NewTaskHandler(mockService, logChan)

	req := httptest.NewRequest(http.MethodPut, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.HandleTasks(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
