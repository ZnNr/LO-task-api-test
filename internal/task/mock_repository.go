package task

import (
	"errors"
	"sync"
)

// MockRepository реализует интерфейс Repository для тестирования
type MockRepository struct {
	mu     sync.Mutex
	tasks  map[int]Task
	nextID int
	err    error // ошибка, которую нужно вернуть (для тестирования ошибок)
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

// SetError позволяет установить ошибку для следующих вызовов
func (m *MockRepository) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

// Reset сбрасывает состояние
func (m *MockRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks = make(map[int]Task)
	m.nextID = 1
	m.err = nil
}

func (m *MockRepository) CreateTask(task *Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.err != nil {
		return m.err
	}

	task.ID = m.nextID
	m.tasks[task.ID] = *task
	m.nextID++
	return nil
}

func (m *MockRepository) GetAllTask(status *Status) ([]Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.err != nil {
		return nil, m.err
	}

	var result []Task
	for _, t := range m.tasks {
		if status == nil || t.Status == *status {
			result = append(result, t)
		}
	}
	return result, nil
}

func (m *MockRepository) GetTaskByID(id int) (*Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.err != nil {
		return nil, m.err
	}

	if t, exists := m.tasks[id]; exists {
		return &t, nil
	}
	return nil, errors.New("task not found")
}
