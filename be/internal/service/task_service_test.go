package service

import (
	"testing"
	"time"

	"gophertodo/backend/internal/domain"
	"gophertodo/backend/internal/repository"
)

func TestAddTask_Success(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task, err := svc.AddTask("完成软件工程作业")
	if err != nil {
		t.Fatalf("AddTask returned unexpected error: %v", err)
	}
	if task == nil {
		t.Fatal("AddTask returned nil task")
	}
	if task.Content != "完成软件工程作业" {
		t.Errorf("Content = %q; want %q", task.Content, "完成软件工程作业")
	}
	if task.Status != domain.Pending {
		t.Errorf("Status = %v; want Pending", task.Status)
	}
	if task.ID <= 0 {
		t.Errorf("ID should be positive, got %d", task.ID)
	}
	if task.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestAddTask_EmptyContent_ReturnsError(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	_, err := svc.AddTask("")
	if err == nil {
		t.Fatal("expected error for empty content, got nil")
	}
}

func TestCompleteTask_Success(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task, _ := svc.AddTask("要完成的任务")

	updated, err := svc.CompleteTask(task.ID)
	if err != nil {
		t.Fatalf("CompleteTask returned unexpected error: %v", err)
	}

	if updated.Status != domain.Completed {
		t.Errorf("Status = %v; want Completed", updated.Status)
	}
	if updated.CompletedAt == nil {
		t.Error("CompletedAt should be set after completion")
	}
}

func TestCompleteTask_NotFound_ReturnsError(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	_, err := svc.CompleteTask(999)
	if err == nil {
		t.Fatal("expected ErrTaskNotFound, got nil")
	}
}

func TestCompleteTask_AlreadyCompleted_ReturnsError(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task, _ := svc.AddTask("已完成的任务")
	_, _ = svc.CompleteTask(task.ID)

	_, err := svc.CompleteTask(task.ID)
	if err == nil {
		t.Fatal("expected error for already-completed task, got nil")
	}
}

func TestListTasks_ReturnsAll(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	_, _ = svc.AddTask("任务1")
	_, _ = svc.AddTask("任务2")

	tasks, err := svc.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks returned error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("ListTasks returned %d tasks; want 2", len(tasks))
	}
}

func TestDeleteTask_Success(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task, _ := svc.AddTask("要删除的任务")
	err := svc.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("DeleteTask returned unexpected error: %v", err)
	}

	_, err = repo.FindByID(task.ID)
	if err == nil {
		t.Fatal("expected error after deletion, got nil")
	}
}

func TestTask_MarkCompleted_SetsStatusAndTime(t *testing.T) {
	task := &domain.Task{
		ID:        1,
		Content:   "领域模型测试",
		Status:    domain.Pending,
		CreatedAt: time.Now().Add(-time.Hour),
	}

	err := task.MarkCompleted(time.Now())
	if err != nil {
		t.Errorf("MarkCompleted should not return error for pending task: %v", err)
	}
	if task.Status != domain.Completed {
		t.Errorf("Status = %v; want Completed", task.Status)
	}
	if task.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	err = task.MarkCompleted(time.Now())
	if err == nil {
		t.Error("MarkCompleted should return error for already-completed task")
	}
}