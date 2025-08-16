// Package handler предоставляет HTTP-хендлеры для управления задачами
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ZnNr/LO-task-api-test/internal/task"
)

// TaskHandler обрабатывает HTTP-запросы для задач
type TaskHandler struct {
	service task.Service
	logChan chan string
}

// NewTaskHandler создаёт новый экземпляр TaskHandler
func NewTaskHandler(service task.Service, logChan chan string) *TaskHandler {
	return &TaskHandler{
		service: service,
		logChan: logChan,
	}
}

// HandleTasks обрабатывает запросы к /tasks (GET и POST)
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetTasks(w, r)
	case http.MethodPost:
		h.handleCreateTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	var statusFilter *task.Status
	if status != "" {
		s := task.Status(status)
		statusFilter = &s
	}

	tasks, err := h.service.GetAllTask(r.Context(), statusFilter)
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tasks)

	h.logChan <- fmt.Sprintf("GET /tasks%s: %d tasks returned", r.URL.RawQuery, len(tasks))
}

func (h *TaskHandler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	var req struct {
		Title  string `json:"title"`
		Status string `json:"status,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	var status task.Status
	if req.Status == "" {
		status = task.StatusPending
	} else {
		status = task.Status(req.Status)
	}

	if err := h.service.CreateTask(r.Context(), req.Title, status); err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	h.logChan <- fmt.Sprintf("POST /tasks: created task with title '%s'", req.Title)
}

// HandleTaskByID обрабатывает запросы к /tasks/{id} (только GET)
func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetTaskByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) handleGetTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	taskObj, err := h.service.GetTaskByID(r.Context(), id)
	if err != nil {
		// Проверяем тип ошибки
		if errors.Is(err, task.ErrTaskNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(taskObj)

	h.logChan <- fmt.Sprintf("GET /tasks/%d: task found", id)
}
