package repository

import (
	"os"
	"path/filepath"
	"testing"

	"gophertodo/backend/internal/domain"
)

// Helper function to create a temporary JSON file for testing
func setupTestRepo(t *testing.T) (*JSONTaskRepository, string) {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tasks.json")
	repo, err := NewJSONTaskRepository(tmpFile)
	if err != nil {
		t.Fatalf("NewJSONTaskRepository returned error: %v", err)
	}
	return repo, tmpFile
}

func TestJSONRepository_Save(t *testing.T) {
	repo, _ := setupTestRepo(t)

	task := &domain.Task{
		Content: "测试任务",
		Status:  domain.StatusPending,
	}

	err := repo.Save(task)
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if task.ID <= 0 {
		t.Errorf("task.ID should be positive after Save, got %d", task.ID)
	}
}

func TestJSONRepository_FindByID(t *testing.T) {
	repo, _ := setupTestRepo(t)

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

func TestJSONRepository_FindByID_NotFound(t *testing.T) {
	repo, _ := setupTestRepo(t)

	_, err := repo.FindByID(999)
	if err == nil {
		t.Fatal("expected ErrTaskNotFound, got nil")
	}
	if err != ErrTaskNotFound {
		t.Errorf("error = %v; want %v", err, ErrTaskNotFound)
	}
}

func TestJSONRepository_FindAll(t *testing.T) {
	repo, _ := setupTestRepo(t)

	for i := 1; i <= 3; i++ {
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

func TestJSONRepository_FindAll_Empty(t *testing.T) {
	repo, _ := setupTestRepo(t)

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("len(tasks) = %d; want 0", len(tasks))
	}
}

func TestJSONRepository_Update(t *testing.T) {
	repo, _ := setupTestRepo(t)

	task := &domain.Task{Content: "原始内容", Status: domain.StatusPending}
	_ = repo.Save(task)

	task.Content = "更新后的内容"
	err := repo.Update(task)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	found, _ := repo.FindByID(task.ID)
	if found.Content != "更新后的内容" {
		t.Errorf("Content = %q; want %q", found.Content, "更新后的内容")
	}
}

func TestJSONRepository_Update_NotFound(t *testing.T) {
	repo, _ := setupTestRepo(t)

	task := &domain.Task{ID: 999, Content: "不存在", Status: domain.StatusPending}
	err := repo.Update(task)
	if err == nil {
		t.Fatal("expected ErrTaskNotFound, got nil")
	}
	if err != ErrTaskNotFound {
		t.Errorf("error = %v; want %v", err, ErrTaskNotFound)
	}
}

func TestJSONRepository_Delete(t *testing.T) {
	repo, _ := setupTestRepo(t)

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

func TestJSONRepository_Delete_NotFound(t *testing.T) {
	repo, _ := setupTestRepo(t)

	err := repo.Delete(999)
	if err == nil {
		t.Fatal("expected ErrTaskNotFound, got nil")
	}
	if err != ErrTaskNotFound {
		t.Errorf("error = %v; want %v", err, ErrTaskNotFound)
	}
}

func TestJSONRepository_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tasks.json")

	// Create repo and save a task
	repo1, err := NewJSONTaskRepository(tmpFile)
	if err != nil {
		t.Fatalf("NewJSONTaskRepository returned error: %v", err)
	}

	task := &domain.Task{Content: "持久化测试", Status: domain.StatusPending}
	_ = repo1.Save(task)

	// Create a new repo from the same file
	repo2, err := NewJSONTaskRepository(tmpFile)
	if err != nil {
		t.Fatalf("NewJSONTaskRepository returned error: %v", err)
	}

	found, err := repo2.FindByID(task.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if found.Content != task.Content {
		t.Errorf("Content = %q; want %q", found.Content, task.Content)
	}
}

func TestJSONRepository_Update_Twice(t *testing.T) {
	repo, _ := setupTestRepo(t)

	task := &domain.Task{Content: "初始", Status: domain.StatusPending}
	_ = repo.Save(task)

	// First update
	task.Content = "第一次更新"
	_ = repo.Update(task)

	// Second update
	task.Content = "第二次更新"
	_ = repo.Update(task)

	found, _ := repo.FindByID(task.ID)
	if found.Content != "第二次更新" {
		t.Errorf("Content = %q; want %q", found.Content, "第二次更新")
	}
}

func TestJSONRepository_Delete_Multiple(t *testing.T) {
	repo, _ := setupTestRepo(t)

	var ids []int
	for i := 0; i < 5; i++ {
		task := &domain.Task{Content: "任务", Status: domain.StatusPending}
		_ = repo.Save(task)
		ids = append(ids, task.ID)
	}

	// Delete first, third, fifth
	_ = repo.Delete(ids[0])
	_ = repo.Delete(ids[2])
	_ = repo.Delete(ids[4])

	tasks, _ := repo.FindAll()
	if len(tasks) != 2 {
		t.Errorf("len(tasks) = %d; want 2", len(tasks))
	}
}

func TestJSONRepository_FindAll_SortedByID(t *testing.T) {
	repo, _ := setupTestRepo(t)

	contents := []string{"zebra", "apple", "banana"}
	for _, content := range contents {
		task := &domain.Task{Content: content, Status: domain.StatusPending}
		_ = repo.Save(task)
	}

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}

	// Tasks should be sorted by ID (ascending), not by content
	if tasks[0].Content != "zebra" {
		t.Errorf("First task content = %q; want %q", tasks[0].Content, "zebra")
	}
}

func TestJSONRepository_ConcurrentSave(t *testing.T) {
	repo, _ := setupTestRepo(t)

	// Save multiple tasks concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			task := &domain.Task{Content: "并发任务", Status: domain.StatusPending}
			_ = repo.Save(task)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}
	if len(tasks) != 10 {
		t.Errorf("len(tasks) = %d; want 10", len(tasks))
	}
}

func TestJSONRepository_Load_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty.json")

	// Create an empty file
	_ = os.WriteFile(tmpFile, []byte(""), 0644)

	repo, err := NewJSONTaskRepository(tmpFile)
	if err != nil {
		t.Fatalf("NewJSONTaskRepository returned error: %v", err)
	}

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("len(tasks) = %d; want 0", len(tasks))
	}
}

func TestJSONRepository_Load_MalformedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "malformed.json")

	// Write malformed JSON
	_ = os.WriteFile(tmpFile, []byte("{invalid json}"), 0644)

	_, err := NewJSONTaskRepository(tmpFile)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

func TestJSONRepository_Save_UpdatesNextID(t *testing.T) {
	repo, _ := setupTestRepo(t)

	task1 := &domain.Task{Content: "任务1", Status: domain.StatusPending}
	task2 := &domain.Task{Content: "任务2", Status: domain.StatusPending}

	_ = repo.Save(task1)
	_ = repo.Save(task2)

	if task2.ID != task1.ID+1 {
		t.Errorf("task2.ID = %d; want %d", task2.ID, task1.ID+1)
	}
}

func TestJSONRepository_FindByID_ReturnsCopy(t *testing.T) {
	repo, _ := setupTestRepo(t)

	task := &domain.Task{Content: "原始内容", Status: domain.StatusPending}
	_ = repo.Save(task)

	// Modify the returned task
	found, _ := repo.FindByID(task.ID)
	found.Content = "修改后"

	// Find again, should be unchanged
	found2, _ := repo.FindByID(task.ID)
	if found2.Content != "原始内容" {
		t.Errorf("Content = %q; want %q", found2.Content, "原始内容")
	}
}

func TestJSONRepository_Update_ReturnsCopy(t *testing.T) {
	repo, _ := setupTestRepo(t)

	task := &domain.Task{Content: "原始", Status: domain.StatusPending}
	_ = repo.Save(task)

	// Get, modify, update
	found, _ := repo.FindByID(task.ID)
	found.Content = "修改"
	_ = repo.Update(found)

	// Original task should be unchanged
	if task.Content != "原始" {
		t.Errorf("Original task.Content = %q; want %q", task.Content, "原始")
	}
}

func TestJSONRepository_Delete_InvalidID(t *testing.T) {
	repo, _ := setupTestRepo(t)

	err := repo.Delete(-1)
	if err == nil {
		t.Fatal("expected error for negative ID, got nil")
	}

	err = repo.Delete(0)
	if err == nil {
		t.Fatal("expected error for zero ID, got nil")
	}
}
