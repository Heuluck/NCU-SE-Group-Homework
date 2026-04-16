package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"gophertodo/backend/internal/domain"
)

type JSONTaskRepository struct {
	mu     sync.RWMutex
	path   string
	nextID int
	tasks  map[int]*domain.Task
}

type taskStore struct {
	NextID int            `json:"next_id"`
	Tasks  []*domain.Task `json:"tasks"`
}

func NewJSONTaskRepository(path string) (*JSONTaskRepository, error) {
	repo := &JSONTaskRepository{
		path:   path,
		nextID: 1,
		tasks:  make(map[int]*domain.Task),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *JSONTaskRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := cloneTask(task)
	copy.ID = r.nextID
	r.nextID++
	r.tasks[copy.ID] = copy
	task.ID = copy.ID

	return r.persistLocked()
}

func (r *JSONTaskRepository) FindByID(id int) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}
	return cloneTask(task), nil
}

func (r *JSONTaskRepository) FindAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, cloneTask(task))
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})
	return tasks, nil
}

func (r *JSONTaskRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[task.ID]; !ok {
		return ErrTaskNotFound
	}
	r.tasks[task.ID] = cloneTask(task)
	return r.persistLocked()
}

func (r *JSONTaskRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[id]; !ok {
		return ErrTaskNotFound
	}
	delete(r.tasks, id)
	return r.persistLocked()
}

func (r *JSONTaskRepository) load() error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return err
	}

	content, err := os.ReadFile(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return r.persistLocked()
	}
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return nil
	}

	var store taskStore
	if err := json.Unmarshal(content, &store); err != nil {
		return err
	}

	maxID := 0
	for _, task := range store.Tasks {
		if task == nil {
			continue
		}
		r.tasks[task.ID] = cloneTask(task)
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	r.nextID = store.NextID
	if r.nextID <= maxID {
		r.nextID = maxID + 1
	}
	if r.nextID < 1 {
		r.nextID = 1
	}
	return nil
}

func (r *JSONTaskRepository) persistLocked() error {
	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, cloneTask(task))
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	store := taskStore{
		NextID: r.nextID,
		Tasks:  tasks,
	}

	content, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, content, 0644)
}

func cloneTask(task *domain.Task) *domain.Task {
	if task == nil {
		return nil
	}
	copy := *task
	if task.CompletedAt != nil {
		completedAt := *task.CompletedAt
		copy.CompletedAt = &completedAt
	}
	return &copy
}
