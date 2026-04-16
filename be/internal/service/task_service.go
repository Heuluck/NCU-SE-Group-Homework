package service

import (
	"time"

	"gophertodo/backend/internal/domain"
	"gophertodo/backend/internal/repository"
)

type TaskAppService struct {
	repo repository.TaskRepository
	now  func() time.Time
}

func NewTaskAppService(repo repository.TaskRepository) *TaskAppService {
	return &TaskAppService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *TaskAppService) AddTask(content string) (*domain.Task, error) {
	task, err := domain.NewTask(content, s.now())
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskAppService) ListTasks() ([]*domain.Task, error) {
	return s.repo.FindAll()
}

func (s *TaskAppService) GetTask(id int) (*domain.Task, error) {
	return s.repo.FindByID(id)
}

func (s *TaskAppService) CompleteTask(id int) (*domain.Task, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := task.MarkCompleted(s.now()); err != nil {
		return nil, err
	}
	if err := s.repo.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskAppService) DeleteTask(id int) error {
	return s.repo.Delete(id)
}
