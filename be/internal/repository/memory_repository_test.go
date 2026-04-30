package repository

import (
	"testing"

	"gophertodo/backend/internal/domain"
)

func TestMemoryRepository_Save(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "测试任务", Status: domain.StatusPending}
	err := repo.Save(task)
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if task.ID <= 0 {
		t.Errorf("task.ID should be positive after Save, got %d", task.ID)
	}
}

func TestMemoryRepository_Save_IncrementsID(t *testing.T) {
	repo := NewMemoryRepository()

	task1 := &domain.Task{Content: "任务1", Status: domain.StatusPending}
	task2 := &domain.Task{Content: "任务2", Status: domain.StatusPending}

	_ = repo.Save(task1)
	_ = repo.Save(task2)

	if task2.ID != task1.ID+1 {
		t.Errorf("task2.ID = %d; want %d", task2.ID, task1.ID+1)
	}
}

func TestMemoryRepository_FindByID(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "测试任务", Status: domain.StatusPending}
	_ = repo.Save(task)

	found, err := repo.FindByID(task.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if found.ID != task.ID {
		t.Errorf("ID = %d; want %d", found.ID, task.ID)
	}
	if found.Content != task.Content {
		t.Errorf("Content = %q; want %q", found.Content, task.Content)
	}
}

func TestMemoryRepository_FindByID_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.FindByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMemoryRepository_FindAll(t *testing.T) {
	repo := NewMemoryRepository()

	for i := 0; i < 3; i++ {
		task := &domain.Task{Content: "任务", Status: domain.StatusPending}
		_ = repo.Save(task)
	}

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("len(tasks) = %d; want 3", len(tasks))
	}
}

func TestMemoryRepository_FindAll_Empty(t *testing.T) {
	repo := NewMemoryRepository()

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("len(tasks) = %d; want 0", len(tasks))
	}
}

func TestMemoryRepository_Update(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "原始内容", Status: domain.StatusPending}
	_ = repo.Save(task)

	task.Content = "更新后"
	err := repo.Update(task)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	found, _ := repo.FindByID(task.ID)
	if found.Content != "更新后" {
		t.Errorf("Content = %q; want %q", found.Content, "更新后")
	}
}

func TestMemoryRepository_Update_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{ID: 999, Content: "不存在", Status: domain.StatusPending}
	err := repo.Update(task)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMemoryRepository_Delete(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "要删除", Status: domain.StatusPending}
	_ = repo.Save(task)

	err := repo.Delete(task.ID)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	_, err = repo.FindByID(task.ID)
	if err == nil {
		t.Fatal("expected error after deletion, got nil")
	}
}

func TestMemoryRepository_Delete_NotFound(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.Delete(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestMemoryRepository_Count(t *testing.T) {
	repo := NewMemoryRepository()

	if repo.Count() != 0 {
		t.Errorf("Initial count = %d; want 0", repo.Count())
	}

	for i := 0; i < 5; i++ {
		task := &domain.Task{Content: "任务", Status: domain.StatusPending}
		_ = repo.Save(task)
	}

	if repo.Count() != 5 {
		t.Errorf("Count = %d; want 5", repo.Count())
	}
}

func TestMemoryRepository_Delete_DecreasesCount(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "任务", Status: domain.StatusPending}
	_ = repo.Save(task)

	if repo.Count() != 1 {
		t.Errorf("Count = %d; want 1", repo.Count())
	}

	_ = repo.Delete(task.ID)

	if repo.Count() != 0 {
		t.Errorf("Count = %d; want 0", repo.Count())
	}
}

func TestMemoryRepository_FindByID_ReturnsCopy(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "原始", Status: domain.StatusPending}
	_ = repo.Save(task)

	found, _ := repo.FindByID(task.ID)
	found.Content = "修改"

	found2, _ := repo.FindByID(task.ID)
	if found2.Content != "原始" {
		t.Errorf("Content = %q; want %q", found2.Content, "原始")
	}
}

func TestMemoryRepository_Update_Twice(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "初始", Status: domain.StatusPending}
	_ = repo.Save(task)

	task.Content = "第一次"
	_ = repo.Update(task)

	task.Content = "第二次"
	_ = repo.Update(task)

	found, _ := repo.FindByID(task.ID)
	if found.Content != "第二次" {
		t.Errorf("Content = %q; want %q", found.Content, "第二次")
	}
}

func TestMemoryRepository_Delete_AfterUpdate(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "任务", Status: domain.StatusPending}
	_ = repo.Save(task)

	task.Content = "已更新"
	_ = repo.Update(task)

	_ = repo.Delete(task.ID)

	if repo.Count() != 0 {
		t.Errorf("Count = %d; want 0", repo.Count())
	}
}

func TestMemoryRepository_Update_AfterDelete(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "任务", Status: domain.StatusPending}
	_ = repo.Save(task)
	_ = repo.Delete(task.ID)

	err := repo.Update(task)
	if err == nil {
		t.Error("expected error when updating deleted task, got nil")
	}
}

func TestMemoryRepository_ConcurrentSave(t *testing.T) {
	repo := NewMemoryRepository()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			task := &domain.Task{Content: "并发", Status: domain.StatusPending}
			_ = repo.Save(task)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	if repo.Count() != 10 {
		t.Errorf("Count = %d; want 10", repo.Count())
	}
}

func TestMemoryRepository_Delete_InvalidID(t *testing.T) {
	repo := NewMemoryRepository()

	err := repo.Delete(-1)
	if err == nil {
		t.Error("expected error for negative ID, got nil")
	}

	err = repo.Delete(0)
	if err == nil {
		t.Error("expected error for zero ID, got nil")
	}
}

func TestMemoryRepository_FindAll_ReturnsCopies(t *testing.T) {
	repo := NewMemoryRepository()

	task := &domain.Task{Content: "原始", Status: domain.StatusPending}
	_ = repo.Save(task)

	tasks, _ := repo.FindAll()
	tasks[0].Content = "修改"

	found, _ := repo.FindByID(task.ID)
	if found.Content != "原始" {
		t.Errorf("Content = %q; want %q", found.Content, "原始")
	}
}
