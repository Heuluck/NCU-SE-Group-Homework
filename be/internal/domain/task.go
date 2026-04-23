package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	Pending   = "pending"
	Completed = "completed"
	StatusPending   = "pending"
	StatusCompleted = "completed"
)

var (
	ErrEmptyContent  = errors.New("task content is required")
	ErrAlreadyDone   = errors.New("task is already completed")
	ErrTaskNotFound = errors.New("task not found")
)

type Task struct {
	ID          int        `json:"id"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func NewTask(content string, now time.Time) (*Task, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ErrEmptyContent
	}

	return &Task{
		Content:   content,
		Status:    StatusPending,
		CreatedAt: now.UTC(),
	}, nil
}

func (t *Task) MarkCompleted(now time.Time) error {
	if t.Status == StatusCompleted {
		return ErrAlreadyDone
	}

	completedAt := now.UTC()
	t.Status = StatusCompleted
	t.CompletedAt = &completedAt
	return nil
}
