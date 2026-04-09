# 详细设计说明书 (DDS) 更新

**项目名称**: GopherTodo CLI
**团队**: 啥队
**文档版本**: v2.0
**更新日期**: 2026-04-09
**关联Sprint**: Sprint 2

---

## 一、文档目的

本文档是GopherTodo项目的详细设计说明书（DDS）更新版本，针对Sprint 2中提取出的核心控制类（TaskAppService、CLIHandler、TaskRepository），归档其内部关键算法的**PAD图**与**N-S盒图**，并规范化跨类调用的方法签名细节。

---

## 二、核心控制类总览

| 类名 | 职责 | 类型 | 所在包 |
|------|------|------|--------|
| `CLIHandler` | 命令解析与输出渲染 | 控制类 | `cmd/` |
| `TaskAppService` | 核心业务逻辑编排 | 应用服务类 | `service/` |
| `Task` | 领域实体与业务行为 | 实体类 | `domain/` |
| `ITaskRepository` | 数据访问抽象接口 | 接口类 | `repository/` |
| `SQLiteRepository` | SQLite持久化实现 | 仓储实现类 | `repository/` |
| `MemoryRepository` | 内存持久化实现(测试用) | 仓储实现类 | `repository/` |

---

## 三、核心算法 PAD 图

### 3.1 TaskAppService.CompleteTask() - PAD图

```
┌──────────────────────────────────────────────────────────────┐
│                    CompleteTask(id int)                      │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  SEQUENCE id                                                 │
│    │                                                         │
│    ▼                                                         │
│  ┌─────────────────┐                                         │
│  │ task, err :=    │  ① 调用仓库查找任务                      │
│  │ repo.FindById   │                                         │
│  │ (id)            │                                         │
│  └───────┬─────────┘                                         │
│          │                                                   │
│          ▼                                                   │
│  ┌─────────────────┐    ┌─────────────────┐                 │
│  │ err != nil ?    │YES │ return          │                 │
│  └───────┬─────────┘    │ ErrTaskNotFound │                 │
│          │NO             └─────────────────┘                 │
│          ▼                                                   │
│  ┌─────────────────┐                                         │
│  │ success :=      │  ② 调用Task领域行为                     │
│  │ task.Mark       │                                         │
│  │ Completed()     │                                         │
│  └───────┬─────────┘                                         │
│          │                                                   │
│          ▼                                                   │
│  ┌─────────────────┐    ┌─────────────────┐                 │
│  │ !success ?      │YES │ return Err      │                 │
│  └───────┬─────────┘    │ AlreadyDone     │                 │
│          │NO             └─────────────────┘                 │
│          ▼                                                   │
│  ┌─────────────────┐                                         │
│  │ return repo.    │  ③ 更新持久化                            │
│  │ Update(task)    │                                         │
│  └───────┬─────────┘                                         │
│          │                                                   │
│          ▼                                                   │
│  ┌─────────────────┐                                         │
│  │ return nil      │  ④ 返回成功                              │
│  └─────────────────┘                                         │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 3.2 Task.MarkCompleted() - PAD图

```
┌──────────────────────────────────────────────────────────────┐
│                   MarkCompleted() bool                       │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────────────────┐                            │
│  │ SELECT                       │                            │
│  │   ┌───────────────────────┐ │                            │
│  │   │ Status == Completed ? │ │  ① 判断当前状态            │
│  │   └───────────┬───────────┘ │                            │
│  │       │YES     │NO          │                            │
│  │       ▼        ▼            │                            │
│  │   ┌─────┐  ┌─────────────┐  │                            │
│  │   │false│  │Status :=    │  │                            │
│  │   │     │  │Completed    │  │                            │
│  │   └──┬──┘  └──────┬──────┘  │                            │
│  │      │            │         │                            │
│  │      │     ┌──────▼──────┐  │                            │
│  │      │     │CompletedAt :=│  │                            │
│  │      │     │time.Now()   │  │                            │
│  │      │     └──────┬──────┘  │                            │
│  │      │            │         │                            │
│  │      │     ┌──────▼──────┐  │                            │
│  │      │     │ return true │  │                            │
│  │      │     └─────────────┘  │                            │
│  │      ▼                      │                            │
│  │   ┌─────────────┐           │                            │
│  │   │ return false│           │                            │
│  │   └─────────────┘           │                            │
│  └─────────────────────────────┘                            │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 3.3 CLIHandler.Handle() - PAD图

```
┌──────────────────────────────────────────────────────────────┐
│                      Handle(args []string)                   │
├──────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ SELECT                                                     │  │
│  │   command := parseCommand(args)                           │  │
│  │   │                                                       │  │
│  │   ▼                                                       │  │
│  │   ┌──────────┬──────────┬──────────┬──────────┐          │  │
│  │   │command== │command== │command== │command== │          │  │
│  │   │"add"?    │"list"?   │"done"?   │"delete"? │          │  │
│  │   └────┬─────┴────┬─────┴────┬─────┴────┬─────┘          │  │
│  │        │          │          │          │                 │  │
│  │        ▼          ▼          ▼          ▼                 │  │
│  │   ┌─────────┐┌─────────┐┌─────────┐┌─────────┐         │  │
│  │   │AddTask  ││ListTask ││DoneTask ││DelTask  │         │  │
│  │   │         ││         ││         ││         │         │  │
│  │   └────┬────┘└────┬────┘└────┬────┘└────┬────┘         │  │
│  │        └──────────┴──────────┴──────────┘                │  │
│  │                        │                                   │  │
│  │                        ▼                                   │  │
│  │              renderOutput(msg)                           │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

---

## 四、N-S 盒图

### 4.1 TaskAppService.AddTask() - N-S盒图

```
┌──────────────────────────────────────────────────────────────────┐
│  PROCEDURE AddTask(content string)                               │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ BEGIN                                                       │  │
│  │   ┌──────────────────────────────────────────────────────┐ │  │
│  │   │ IF content == "" THEN                                │ │  │
│  │   │   RETURN ErrEmptyContent                             │ │  │
│  │   │ END IF                                               │ │  │
│  │   └──────────────────────────────────────────────────────┘ │  │
│  │   │                                                       │ │  │
│  │   ├──────────────────────────────────────────────────────┤ │  │
│  │   │ task := Task{                                        │ │  │
│  │   │         Content:   content,                          │ │  │
│  │   │         Status:    Pending,                          │ │  │
│  │   │         CreatedAt: time.Now()                         │ │  │
│  │   │       }                                              │ │  │
│  │   └──────────────────────────────────────────────────────┘ │  │
│  │   │                                                       │ │  │
│  │   ├──────────────────────────────────────────────────────┤ │  │
│  │   │ IF err := repo.Save(&task); err != nil THEN          │ │  │
│  │   │   RETURN err                                         │ │  │
│  │   │ END IF                                               │ │  │
│  │   └──────────────────────────────────────────────────────┘ │  │
│  │   │                                                       │ │  │
│  │   ├──────────────────────────────────────────────────────┤ │  │
│  │   │ RETURN nil                                           │ │  │
│  │   └──────────────────────────────────────────────────────┘ │  │
│  │ END                                                         │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 4.2 TaskAppService.ListTasks() - N-S盒图

```
┌──────────────────────────────────────────────────────────────────┐
│  PROCEDURE ListTasks() []*Task                                   │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ BEGIN                                                       │  │
│  │   ┌──────────────────────────────────────────────────────┐ │  │
│  │   │ tasks, err := repo.FindAll()                         │ │  │
│  │   │ IF err != nil THEN                                   │ │  │
│  │   │   RETURN nil, err                                    │ │  │
│  │   │ END IF                                               │ │  │
│  │   └──────────────────────────────────────────────────────┘ │  │
│  │   │                                                       │ │  │
│  │   ├──────────────────────────────────────────────────────┤ │  │
│  │   │ filtered := make([]*Task, 0)                         │ │  │
│  │   │ FOR EACH task IN tasks DO                             │ │  │
│  │   │   ┌────────────────────────────────────────────────┐ │ │  │
│  │   │   │ IF task.Status == Pending THEN                  │ │ │  │
│  │   │   │   filtered = append(filtered, task)            │ │ │  │
│  │   │   │ END IF                                          │ │ │  │
│  │   │   └────────────────────────────────────────────────┘ │ │  │
│  │   │ END FOR                                              │ │  │
│  │   └──────────────────────────────────────────────────────┘ │  │
│  │   │                                                       │ │  │
│  │   ├──────────────────────────────────────────────────────┤ │  │
│  │   │ RETURN filtered                                      │ │  │
│  │   └──────────────────────────────────────────────────────┘ │  │
│  │ END                                                         │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 五、跨类调用方法签名规范

### 5.1 CLIHandler → TaskAppService

| 调用方法 | 签名 | 返回值 | 说明 |
|---------|------|--------|------|
| 添加任务 | `appService.AddTask(content string)` | `(*Task, error)` | 创建新任务 |
| 查看列表 | `appService.ListTasks()` | `([]*Task, error)` | 获取所有未完成任务 |
| 完成任务 | `appService.CompleteTask(id int)` | `error` | 根据ID标记完成 |
| 删除任务 | `appService.DeleteTask(id int)` | `error` | 根据ID删除任务 |

### 5.2 TaskAppService → ITaskRepository

| 调用方法 | 签名 | 返回值 | 说明 |
|---------|------|--------|------|
| 保存 | `repo.Save(task *Task)` | `error` | 持久化新任务 |
| 查询单个 | `repo.FindById(id int)` | `(*Task, error)` | 根据ID查询 |
| 查询全部 | `repo.FindAll()` | `([]*Task, error)` | 查询所有任务 |
| 更新 | `repo.Update(task *Task)` | `error` | 更新任务状态 |
| 删除 | `repo.Delete(id int)` | `error` | 删除指定任务 |

### 5.3 完整调用链示例

```
用户输入: todo done 1
        │
        ▼
CLIHandler.Handle(["done", "1"])
        │
        ▼
parseCommand() → command="done", args=["1"]
        │
        ▼
TaskAppService.CompleteTask(1)
        │
        ├─→ ITaskRepository.FindById(1)
        │         │
        │         ▼
        │   SQLiteRepository.FindById(1)
        │         │
        │         ▼
        │   SELECT * FROM tasks WHERE id=1
        │
        ├─→ Task.MarkCompleted()
        │
        └─→ ITaskRepository.Update(task)
                  │
                  ▼
            SQLiteRepository.Update(task)
                  │
                  ▼
            UPDATE tasks SET status='done' WHERE id=1
```

---

## 六、接口契约定义

### 6.1 ITaskRepository 接口

```go
package repository

// ITaskRepository 定义数据访问抽象接口
// 实现类: SQLiteRepository, MemoryRepository
type ITaskRepository interface {
    // Save 持久化新任务
    Save(task *Task) error

    // FindById 根据ID查找任务
    FindById(id int) (*Task, error)

    // FindAll 返回所有任务
    FindAll() ([]*Task, error)

    // Update 更新已存在的任务
    Update(task *Task) error

    // Delete 根据ID删除任务
    Delete(id int) error
}
```

### 6.2 TaskAppService 结构

```go
package service

// TaskAppService 核心应用服务
// 依赖倒置: 通过接口ITaskRepository编程，不直接依赖具体实现
type TaskAppService struct {
    repo repository.ITaskRepository
}

// NewTaskAppService 构造函数注入依赖
func NewTaskAppService(repo repository.ITaskRepository) *TaskAppService {
    return &TaskAppService{repo: repo}
}
```

---

## 七、数据库表结构

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    content     TEXT    NOT NULL,
    status      TEXT    NOT NULL DEFAULT 'pending',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME
);

CREATE INDEX idx_tasks_status ON tasks(status);
```

---

**文档状态**: 已审核
**下次更新**: Sprint 3 结束时
