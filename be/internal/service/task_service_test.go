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

// ============================================================================
// Prompt 演化实验案例 2：Service 层 CompleteTask 错误处理
// ============================================================================
//
// 【初始 Prompt】
// "测试 CompleteTask 方法"
//
// 【AI 输出的问题】
// - 缺少对 ErrTaskNotFound 的独立测试用例
// - 缺少对 ErrAlreadyDone 的独立测试用例
// - 测试用例混杂在一起，难以快速定位问题
//
// 【改进后的 Prompt（角色扮演 + CoT）】
// "你是一个 Go 测试专家，遵循 AAA 模式（Arrange-Act-Assert）。
// 请为 TaskAppService.CompleteTask 方法编写独立的错误路径测试。
// 思考过程（CoT）：
// 1. CompleteTask 可能返回的错误：
//    - ErrTaskNotFound：当任务不存在时
//    - ErrAlreadyDone：当任务已完成时
// 2. 每个错误路径需要独立测试
// 3. 正常路径（成功完成）也需要测试
// 请生成以下测试函数：
// - TestCompleteTask_NotFound_ReturnsError
// - TestCompleteTask_AlreadyCompleted_ReturnsError
// - TestCompleteTask_Success
//
// 【最终可用测试代码】见下方：
// ============================================================================

func TestGetTask_Success(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	created, _ := svc.AddTask("要获取的任务")

	got, err := svc.GetTask(created.ID)
	if err != nil {
		t.Fatalf("GetTask returned unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %d; want %d", got.ID, created.ID)
	}
	if got.Content != created.Content {
		t.Errorf("Content = %q; want %q", got.Content, created.Content)
	}
	if got.Status != created.Status {
		t.Errorf("Status = %v; want %v", got.Status, created.Status)
	}
}

func TestGetTask_NotFound_ReturnsError(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	_, err := svc.GetTask(999)
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}
}

func TestDeleteTask_NotFound_ReturnsError(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	err := svc.DeleteTask(999)
	if err == nil {
		t.Fatal("expected error for non-existent task, got nil")
	}
}

func TestDeleteTask_AfterComplete(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task, _ := svc.AddTask("先完成后删除")
	_, _ = svc.CompleteTask(task.ID)

	err := svc.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("DeleteTask returned unexpected error: %v", err)
	}

	_, err = svc.GetTask(task.ID)
	if err == nil {
		t.Fatal("expected error after deletion, got nil")
	}
}

func TestListTasks_Empty(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	tasks, err := svc.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks returned error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("ListTasks returned %d tasks; want 0", len(tasks))
	}
}

func TestListTasks_AfterDelete(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task1, _ := svc.AddTask("任务1")
	task2, _ := svc.AddTask("任务2")

	_ = svc.DeleteTask(task1.ID)

	tasks, err := svc.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks returned error: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("ListTasks returned %d tasks; want 1", len(tasks))
	}
	if tasks[0].ID != task2.ID {
		t.Errorf("Remaining task ID = %d; want %d", tasks[0].ID, task2.ID)
	}
}

func TestAddTask_TrimsContent(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task, err := svc.AddTask("  任务内容  ")
	if err != nil {
		t.Fatalf("AddTask returned error: %v", err)
	}
	if task.Content != "任务内容" {
		t.Errorf("Content = %q; want %q", task.Content, "任务内容")
	}
}

func TestAddTask_MultipleTasks(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	contents := []string{"任务1", "任务2", "任务3"}
	var ids []int

	for _, content := range contents {
		task, err := svc.AddTask(content)
		if err != nil {
			t.Fatalf("AddTask(%q) returned error: %v", content, err)
		}
		ids = append(ids, task.ID)
	}

	// Verify all IDs are unique
	seen := make(map[int]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		seen[id] = true
	}

	// Verify all tasks are retrievable
	for i, id := range ids {
		task, err := svc.GetTask(id)
		if err != nil {
			t.Fatalf("GetTask(%d) returned error: %v", id, err)
		}
		if task.Content != contents[i] {
			t.Errorf("Content = %q; want %q", task.Content, contents[i])
		}
	}
}

func TestCompleteTask_SetsCompletedAt(t *testing.T) {
	repo := repository.NewMemoryRepository()
	svc := NewTaskAppService(repo)

	task, _ := svc.AddTask("任务")
	before := time.Now()

	completed, err := svc.CompleteTask(task.ID)
	if err != nil {
		t.Fatalf("CompleteTask returned error: %v", err)
	}

	if completed.CompletedAt == nil {
		t.Fatal("CompletedAt should be set")
	}
	if completed.CompletedAt.Before(before) {
		t.Error("CompletedAt should be after or equal to the time before completion")
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