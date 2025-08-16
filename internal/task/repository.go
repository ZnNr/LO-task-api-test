package task

import "sync"

type Repository struct {
	tasks  map[int]Task
	nextID int
	mu     sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (r *Repository) Create(t *Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t.ID = r.nextID
	r.tasks[t.ID] = *t
	r.nextID++
}

func (r *Repository) GetAll(status *Status) []Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Task
	for _, t := range r.tasks {
		if status == nil || t.Status == *status {
			result = append(result, t)
		}
	}
	return result
}

func (r *Repository) GetByID(id int) (*Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	return &t, ok
}
