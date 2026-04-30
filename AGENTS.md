# GopherTodo — AI 编程助手规范

## 1. 项目架构概述

GopherTodo 是一款**全栈待办事项管理系统**，采用 **Go 后端 + React 前端** 的前后端分离架构。后端遵循 DDD 分层架构：Domain（领域实体）→ Service（应用服务）→ Repository（仓储接口）→ Infrastructure（JSON/Memory 实现）。

**技术栈**：
- 后端：Go 1.23+, 标准库 net/http
- 前端：React 18, TypeScript, Ant Design, Vite
- 测试：Go testing, Vitest + @testing-library/react

---

## 2. 目录结构说明

```
be/                          # 后端（Go）
├── cmd/server/main.go       # HTTP API 服务器入口
├── internal/
│   ├── domain/task.go       # ★ 核心领域实体：Task 实体 + 业务规则
│   ├── service/task_service.go  # 应用服务：AddTask/ListTasks/CompleteTask/DeleteTask
│   ├── repository/          # 仓储层（接口 + 实现）
│   │   ├── task_repository.go   # 接口定义
│   │   ├── json_repository.go   # JSON 文件持久化实现
│   │   └── memory_repository.go  # 内存实现（测试专用）
│   └── httpapi/server.go    # HTTP 处理函数
├── data/tasks.json          # 开发环境 JSON 数据文件
└── go.mod

fe/                          # 前端（React）
├── app/
│   ├── lib/api.ts           # Axios HTTP 客户端
│   ├── components/          # React 组件
│   └── routes/home.tsx      # 首页路由
├── tests/                   # 前端测试
└── openapi.yaml             # OpenAPI 规范
```

---

## 3. 核心模块职责

| 模块 | 文件 | 职责 | 关键规则 |
|------|------|------|---------|
| **Domain** | `domain/task.go` | Task 实体 + 业务规则 | `NewTask`: 空内容返回 `ErrEmptyContent`；`MarkCompleted`: 已完成返回 `ErrAlreadyDone` |
| **Service** | `service/task_service.go` | 业务流程编排 | 依赖 Repository 接口，核心方法：AddTask/ListTasks/GetTask/CompleteTask/DeleteTask |
| **Repository** | `repository/*.go` | 持久化抽象 | MemoryRepository 用于测试，JSONRepository 用于生产 |
| **HTTP API** | `httpapi/server.go` | REST 接口 | `/tasks` GET/POST, `/tasks/{id}` GET/DELETE, `/tasks/{id}/complete` POST |

---

## 4. 编码规范约束

### Go 规范
- **错误处理**：使用 sentinel errors（`ErrTaskNotFound`, `ErrEmptyContent`, `ErrAlreadyDone`），不接受裸 `error`
- **命名**：驼峰式（`NewTaskAppService`），私有函数前导下划线仅用于避免冲突
- **测试**：每个导出函数需有 `_test.go`，使用 AAA 模式（Arrange-Act-Assert）
- **覆盖率**：domain/service/repository ≥ 80%

### TypeScript 规范
- **类型**：`interface Task { id: number; content: string; status: "pending"|"completed"; ... }`
- **API 封装**：`fe/app/lib/api.ts` 中统一封装，错误响应 `{ status: number; data: { code: string; message: string } }`

### 禁止事项
- ❌ 禁止在 `domain/task.go` 外直接修改 Task 实体的 Status/CompletedAt
- ❌ 禁止在 Service 层直接操作 Repository 实现类
- ❌ 禁止绕过 `TaskRepository` 接口进行硬编码依赖
- ❌ 禁止在 `httpapi` 层处理业务逻辑（如校验内容非空）
- ❌ 禁止提交未通过 `go test ./...` 的代码

---

## 5. 禁止操作清单

1. **禁止修改 Domain 业务规则**：Task 实体的 `NewTask`/`MarkCompleted` 是核心领域逻辑，AI 不得擅自修改
2. **禁止删除或削弱错误处理**：所有 `ErrTaskNotFound`/`ErrEmptyContent`/`ErrAlreadyDone` 必须保留
3. **禁止绕过 CI 测试**：PR 必须通过 `go test ./...` 和 `pnpm test`
4. **禁止硬编码敏感信息**：数据库路径、API 地址必须通过环境变量或配置文件
5. **禁止降低测试覆盖率**：新增代码必须附带测试，覆盖率不得低于 80%
