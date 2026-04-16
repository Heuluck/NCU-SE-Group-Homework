package repository

import (
	"errors"

	"gophertodo/backend/internal/domain"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Save(task *domain.Task) error
	FindByID(id int) (*domain.Task, error)
	FindAll() ([]*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id int) error
}
