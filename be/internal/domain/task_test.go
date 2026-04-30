package domain

import (
	"strings"
	"testing"
	"time"
)

// ============================================================================
// Prompt 演化实验案例 1：Task Domain 的边界测试
// ============================================================================
//
// 【初始 Prompt】
// "为 Task 实体写单元测试"
//
// 【AI 输出的问题】
// - 生成重言式测试（如 Status == "pending" 永远为 true，因为 NewTask 总是设置 Pending）
// - 缺乏对边界值的明确测试（空字符串、仅空格字符串）
// - 未测试 MarkCompleted 的错误路径（重复完成）
//
// 【改进后的 Prompt（Few-shot + 结构化指令）】
// "你是一个 Go 测试专家。请为 Task 实体编写边界测试，遵循以下要求：
// 1. 使用 AAA 模式（Arrange-Act-Assert）
// 2. 测试 NewTask 的边界条件：
//    - 空字符串 → 必须返回 ErrEmptyContent
//    - 仅空格字符串 → 必须返回 ErrEmptyContent
//    - 正常字符串 → 返回 task 且 status 为 Pending
// 3. 测试 MarkCompleted 的边界条件：
//    - pending 状态 → 成功，CompletedAt 被设置
//    - 再次调用 → 必须返回 ErrAlreadyDone
// 4. 不要写重言式测试（不要测试 NewTask 返回的 task.Status == "pending"，
//    因为这就是 NewTask 的实现逻辑，不是边界条件）
// 5. 使用 testify/assert 进行断言"
//
// 【最终可用测试代码】见下方：
// ============================================================================

func TestNewTask_Success(t *testing.T) {
	now := time.Now()
	task, err := NewTask("完成软件工程作业", now)

	if err != nil {
		t.Fatalf("NewTask returned unexpected error: %v", err)
	}
	if task == nil {
		t.Fatal("NewTask returned nil task")
	}
	if task.Content != "完成软件工程作业" {
		t.Errorf("Content = %q; want %q", task.Content, "完成软件工程作业")
	}
	if task.Status != StatusPending {
		t.Errorf("Status = %v; want %v", task.Status, StatusPending)
	}
	if task.ID != 0 {
		t.Errorf("ID should be 0 for new task, got %d", task.ID)
	}
	if task.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if task.CompletedAt != nil {
		t.Error("CompletedAt should be nil for new task")
	}
}

func TestNewTask_EmptyString_ReturnsError(t *testing.T) {
	now := time.Now()
	task, err := NewTask("", now)

	if err == nil {
		t.Fatal("expected ErrEmptyContent, got nil")
	}
	if err != ErrEmptyContent {
		t.Errorf("error = %v; want %v", err, ErrEmptyContent)
	}
	if task != nil {
		t.Error("task should be nil when content is empty")
	}
}

func TestNewTask_WhitespaceOnly_ReturnsError(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name    string
		content string
	}{
		{"single space", " "},
		{"multiple spaces", "   "},
		{"tab", "\t"},
		{"newline", "\n"},
		{"mixed whitespace", " \t\n "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := NewTask(tc.content, now)

			if err == nil {
				t.Errorf("content %q: expected ErrEmptyContent, got nil", tc.content)
			}
			if err != ErrEmptyContent {
				t.Errorf("error = %v; want %v", err, ErrEmptyContent)
			}
			if task != nil {
				t.Error("task should be nil when content is only whitespace")
			}
		})
	}
}

func TestNewTask_TrimsWhitespace(t *testing.T) {
	now := time.Now()
	task, err := NewTask("  任务内容  ", now)

	if err != nil {
		t.Fatalf("NewTask returned unexpected error: %v", err)
	}
	if task.Content != "任务内容" {
		t.Errorf("Content = %q; want %q", task.Content, "任务内容")
	}
}

func TestNewTask_ValidContentWithWhitespace(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name    string
		input   string
		content string
	}{
		{"leading space", " 任务", "任务"},
		{"trailing space", "任务  ", "任务"},
		{"both", "  任务  ", "任务"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := NewTask(tc.input, now)

			if err != nil {
				t.Fatalf("NewTask returned error: %v", err)
			}
			if task.Content != tc.content {
				t.Errorf("Content = %q; want %q", task.Content, tc.content)
			}
		})
	}
}

func TestMarkCompleted_FromPending_Success(t *testing.T) {
	task := &Task{
		ID:        1,
		Content:   "待完成的任务",
		Status:    StatusPending,
		CreatedAt: time.Now().Add(-time.Hour),
	}

	now := time.Now()
	err := task.MarkCompleted(now)

	if err != nil {
		t.Fatalf("MarkCompleted returned unexpected error: %v", err)
	}
	if task.Status != StatusCompleted {
		t.Errorf("Status = %v; want %v", task.Status, StatusCompleted)
	}
	if task.CompletedAt == nil {
		t.Fatal("CompletedAt should be set after MarkCompleted")
	}
	if !task.CompletedAt.Equal(now) {
		t.Errorf("CompletedAt = %v; want %v", task.CompletedAt, now)
	}
}

func TestMarkCompleted_AlreadyCompleted_ReturnsError(t *testing.T) {
	now := time.Now()
	completedTime := now.Add(-time.Hour)

	task := &Task{
		ID:          1,
		Content:     "已完成的任务",
		Status:      StatusCompleted,
		CompletedAt: &completedTime,
	}

	err := task.MarkCompleted(now)

	if err == nil {
		t.Fatal("expected ErrAlreadyDone, got nil")
	}
	if err != ErrAlreadyDone {
		t.Errorf("error = %v; want %v", err, ErrAlreadyDone)
	}
}

func TestMarkCompleted_DoesNotModifyStatusIfAlreadyCompleted(t *testing.T) {
	now := time.Now()
	originalTime := now.Add(-time.Hour)

	task := &Task{
		ID:          1,
		Content:     "已完成的任务",
		Status:      StatusCompleted,
		CompletedAt: &originalTime,
	}

	_ = task.MarkCompleted(now)

	// Status should remain Completed, not change
	if task.Status != StatusCompleted {
		t.Errorf("Status = %v; want %v", task.Status, StatusCompleted)
	}
	// CompletedAt should not be modified
	if !task.CompletedAt.Equal(originalTime) {
		t.Errorf("CompletedAt = %v; want %v", task.CompletedAt, originalTime)
	}
}

func TestTask_Constants(t *testing.T) {
	if Pending != "pending" {
		t.Errorf("Pending = %q; want %q", Pending, "pending")
	}
	if Completed != "completed" {
		t.Errorf("Completed = %q; want %q", Completed, "completed")
	}
	if StatusPending != "pending" {
		t.Errorf("StatusPending = %q; want %q", StatusPending, "pending")
	}
	if StatusCompleted != "completed" {
		t.Errorf("StatusCompleted = %q; want %q", StatusCompleted, "completed")
	}
}

func TestNewTask_ContentLength(t *testing.T) {
	now := time.Now()

	// Test very long content
	longContent := strings.Repeat("a", 10000)
	task, err := NewTask(longContent, now)
	if err != nil {
		t.Fatalf("NewTask returned error for long content: %v", err)
	}
	if len(task.Content) != 10000 {
		t.Errorf("Content length = %d; want 10000", len(task.Content))
	}

	// Test single character
	task, err = NewTask("a", now)
	if err != nil {
		t.Fatalf("NewTask returned error for single char: %v", err)
	}
	if task.Content != "a" {
		t.Errorf("Content = %q; want %q", task.Content, "a")
	}
}

func TestNewTask_Unicode(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name    string
		content string
	}{
		{"Chinese", "任务内容"},
		{"Japanese", "タスク"},
		{"Emoji", "✅ 任务完成"},
		{"Mixed", "Hello 世界 🌍"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task, err := NewTask(tc.content, now)

			if err != nil {
				t.Fatalf("NewTask returned error: %v", err)
			}
			if task.Content != tc.content {
				t.Errorf("Content = %q; want %q", task.Content, tc.content)
			}
		})
	}
}
