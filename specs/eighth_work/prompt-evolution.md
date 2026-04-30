# AI 辅助测试的 Prompt 演化实验

## 概述

本文档记录了使用 AI 辅助编写核心 API 模块单元测试的 Prompt 演化过程。通过结构化指令、Few-shot、CoT 等技术，将低质量的初始 Prompt 优化为可用的测试代码。

**核心模块**：Task Domain / Service / Repository（分层架构）

**目标覆盖率**：≥80%

---

## 案例 1：Task Domain 的边界测试

### 初始 Prompt

```
为 Task 实体写单元测试
```

### AI 输出的问题

1. **重言式测试**：生成的测试代码 `if task.Status == "pending"` 永远为 true，因为 `NewTask` 的实现逻辑本身就是设置 Status 为 "pending"。这不是边界条件测试，而是实现细节验证。

2. **缺乏边界覆盖**：
   - 未测试空字符串输入
   - 未测试仅包含空白字符的字符串
   - 未测试 TrimSpace 处理

3. **缺少错误路径测试**：`MarkCompleted` 的重复完成错误路径未被独立测试。

4. **断言方式不规范**：使用 `t.Error` 而非 `t.Fatal`，导致测试不会在第一次失败时停止。

### 改进后的 Prompt

```
你是一个 Go 测试专家。请为 Task 实体编写边界测试，遵循以下要求：

1. 使用 AAA 模式（Arrange-Act-Assert）
2. 测试 NewTask 的边界条件：
   - 空字符串 → 必须返回 ErrEmptyContent
   - 仅空格字符串 → 必须返回 ErrEmptyContent（测试 " "、"\t"、"\n" 等）
   - 正常字符串 → 返回 task 且 status 为 StatusPending
3. 测试 MarkCompleted 的边界条件：
   - pending 状态 → 成功，CompletedAt 被设置
   - 再次调用 → 必须返回 ErrAlreadyDone
4. 不要写重言式测试（不要测试 NewTask 返回的 task.Status == StatusPending，
   因为这就是 NewTask 的实现逻辑，不是边界条件）
5. 使用 t.Fatal/Error 进行断言，确保失败时能快速定位
```

### 最终可用测试代码

```go
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
```

### 演化总结

| 阶段 | Prompt 特征 | 问题 |
|------|------------|------|
| 初始 | 简单指令 | 生成重言式测试、缺乏边界覆盖 |
| 改进 | Few-shot + 结构化指令 | 明确边界值、禁止重言式测试、AAA 模式 |

---

## 案例 2：Service 层 CompleteTask 错误处理

### 初始 Prompt

```
测试 CompleteTask 方法
```

### AI 输出的问题

1. **错误路径覆盖不足**：只测试了"成功完成"的情况，未对 `ErrTaskNotFound` 和 `ErrAlreadyDone` 进行独立测试。

2. **测试用例混杂**：将正常路径和错误路径混在一个测试中，难以快速定位问题。

3. **缺乏边界思考**：未考虑：
   - 不存在的任务 ID
   - 重复完成同一任务
   - 并发完成同一任务

### 改进后的 Prompt（角色扮演 + CoT）

```
你是一个 Go 测试专家，遵循 AAA 模式（Arrange-Act-Assert）。
请为 TaskAppService.CompleteTask 方法编写独立的错误路径测试。

思考过程（CoT）：
1. CompleteTask 可能返回的错误：
   - ErrTaskNotFound：当任务不存在时
   - ErrAlreadyDone：当任务已完成时
2. 每个错误路径需要独立测试函数
3. 正常路径（成功完成）也需要测试
4. 额外考虑：CompletedAt 时间戳是否正确设置

请生成以下测试函数：
- TestCompleteTask_NotFound_ReturnsError
- TestCompleteTask_AlreadyCompleted_ReturnsError
- TestCompleteTask_Success
- TestCompleteTask_SetsCompletedAt（验证时间戳）
```

### 最终可用测试代码

```go
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
```

### 演化总结

| 阶段 | Prompt 特征 | 问题 |
|------|------------|------|
| 初始 | 简单指令 | 缺少独立错误路径测试 |
| 改进 | 角色扮演 + CoT | 明确思考过程、分离测试用例 |

---

## 测试覆盖率报告

### 运行命令

```bash
cd be
go test -coverprofile=coverage.out ./internal/domain/ ./internal/service/ ./internal/repository/
go tool cover -func=coverage.out
```

### 覆盖率结果

| 模块 | 覆盖率 |
|------|--------|
| `internal/domain/` | **100.0%** |
| `internal/service/` | **88.9%** |
| `internal/repository/` | **92.4%** |

### 详细分析

**Domain (100%)**：
- `NewTask`: 100%
- `MarkCompleted`: 100%

**Service (88.9%)**：
- `AddTask`: 83.3%（边界分支未完全覆盖）
- `ListTasks`: 100%
- `GetTask`: 100%
- `CompleteTask`: 87.5%（时间比较分支）
- `DeleteTask`: 100%

**Repository (92.4%)**：
- `Save/FindByID/FindAll/Update/Delete`: 100%
- `load/persistLocked`: 80-90%（部分错误路径）
- `cloneTask`: 57.1%（nil 检查分支）

---

## Prompt 优化技巧总结

### 1. Few-shot（少样本学习）

提供具体示例，明确告诉 AI 什么样的测试是合格的：

```markdown
示例（正确的边界测试）：
func TestNewTask_EmptyString_ReturnsError(t *testing.T) {
    // Arrange
    now := time.Now()
    // Act
    task, err := NewTask("", now)
    // Assert
    if err != ErrEmptyContent { ... }
}
```

### 2. 结构化指令

将需求拆分为明确的编号列表：

```markdown
1. 测试空字符串 → 返回 ErrEmptyContent
2. 测试仅空格 → 返回 ErrEmptyContent
3. 测试正常字符串 → 返回 task 且 Status 为 Pending
```

### 3. CoT（思维链）

要求 AI 先思考再输出，引导其覆盖所有边界：

```markdown
思考过程：
1. 可能返回的错误：...
2. 每个错误需要独立测试
3. 正常路径也需要验证
```

### 4. 角色扮演

指定 AI 的身份，提高输出质量：

```markdown
你是一个 Go 测试专家，遵循 AAA 模式（Arrange-Act-Assert）
```

### 5. 禁止事项

明确告诉 AI 不要做什么：

```markdown
不要写重言式测试（不要测试 NewTask 返回的 task.Status == "pending"，
因为这就是 NewTask 的实现逻辑，不是边界条件）
```

---

## 结论

通过 Prompt 演化，我们成功将初始的低质量测试 Prompt 优化为可用的测试代码，覆盖率达到：

- **Domain**: 100% ✓
- **Service**: 88.9% ✓
- **Repository**: 92.4% ✓

关键改进在于：
1. 从简单指令 → 结构化指令
2. 添加边界值明确说明
3. 引入 Few-shot 示例
4. 使用 CoT 引导完整思考
5. 明确禁止重言式测试
