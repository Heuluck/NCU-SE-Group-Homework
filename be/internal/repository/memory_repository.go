package repository

import (
	"fmt"
	"sync"

	"gophertodo/backend/internal/domain"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	tasks  map[int]*domain.Task
	nextID int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks:  make(map[int]*domain.Task),
		nextID: 0,
	}
}

func (r *MemoryRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	task.ID = r.nextID
	copied := *task
	r.tasks[task.ID] = &copied
	return nil
}

func (r *MemoryRepository) FindByID(id int) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task with id %d not found: %w", id, domain.ErrTaskNotFound)
	}
	copied := *task
	return &copied, nil
}

func (r *MemoryRepository) FindAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		copied := *task
		result = append(result, &copied)
	}
	return result, nil
}

func (r *MemoryRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[task.ID]; !ok {
		return fmt.Errorf("task with id %d not found: %w", task.ID, domain.ErrTaskNotFound)
	}
	copied := *task
	r.tasks[task.ID] = &copied
	return nil
}

func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return fmt.Errorf("task with id %d not found: %w", id, domain.ErrTaskNotFound)
	}
	delete(r.tasks, id)
	return nil
}

func (r *MemoryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tasks)
}