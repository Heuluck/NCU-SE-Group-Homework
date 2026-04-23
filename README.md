# GopherTodo

> AI赋能，啥都能成！ — **啥队** · 南昌大学计科软件工程小组

[![CI — Test & Build](https://github.com/Heuluck/NCU-SE-Group-Homework/actions/workflows/ci-test.yml/badge.svg)](https://github.com/Heuluck/NCU-SE-Group-Homework/actions)
[![Deploy FE Page](https://github.com/Heuluck/NCU-SE-Group-Homework/actions/workflows/fe.yml/badge.svg)](https://github.com/Heuluck/NCU-SE-Group-Homework/actions)

GopherTodo 是一款**全栈待办事项管理系统**，提供命令行（CLI）和 Web 前端两种使用方式，采用 Golang 后端 + React 前端的前后端分离架构。

---

## 目录

- [团队信息](#团队信息)
- [系统架构](#系统架构)
- [快速上手](#快速上手)
  - [环境要求](#环境要求)
  - [后端启动](#后端启动)
  - [前端启动](#前端启动)
  - [使用 Mock 开发](#使用-mock-开发)
- [核心模块职责说明](#核心模块职责说明)
- [API 接口](#api-接口)
- [测试](#测试)
- [项目结构](#项目结构)

---

## 团队信息

| 姓名   | 角色                         | 主要职责                                               |
| ------ | ---------------------------- | ------------------------------------------------------ |
| 陆永祥 | 产品负责人，开发团队（前端） | 产品方向规划、优先级排期；前端框架选型、页面开发、交付 |
| 邓枭   | Scrum Master                 | 主持敏捷会议、落地 Scrum 流程、促进团队协作            |
| 马骏   | 开发团队（后端）             | 后端框架选型、领域模型设计、API 开发、交付             |

---

## 系统架构

### 总体架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                          用户界面层                                 │
│                                                                     │
│   ┌──────────────────────┐        ┌──────────────────────────────┐  │
│   │  CLI（命令行界面）    │        │    Web 前端（React SPA）     │  │
│   │  cobra 子命令路由    │        │    Ant Design + React Router │  │
│   │  todo add/list/done  │        │    fe/app/components/        │  │
│   └──────────┬───────────┘        └──────────────┬───────────────┘  │
└──────────────┼────────────────────────────────────┼─────────────────┘
               │ 方法调用                           │ HTTP / REST API
┌──────────────┼────────────────────────────────────┼─────────────────┐
│              ▼           应用服务层                ▼                │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │                    TaskAppService                           │   │
│   │  AddTask / ListTasks / CompleteTask / DeleteTask / GetTask  │   │
│   └──────────────────────────────┬──────────────────────────────┘   │
└─────────────────────────────────-┼──────────────────────────────────┘
                                   │ 接口调用（依赖倒置）
┌──────────────────────────────────┼──────────────────────────────────┐
│                领域层             │                                  │
│   ┌──────────────────┐           │                                  │
│   │    Task 实体      │    ┌──────┴──────────────────────────────┐  │
│   │  Content / Status│    │       <<interface>>                 │  │
│   │  CreatedAt       │    │       ITaskRepository               │  │
│   │  CompletedAt     │    │  Save / FindById / FindAll          │  │
│   │  MarkCompleted() │    │  Update / Delete                    │  │
│   └──────────────────┘    └──────┬──────────────────────────────┘  │
└─────────────────────────────────-┼──────────────────────────────────┘
                                   │ 多态实现
┌──────────────────────────────────┼──────────────────────────────────┐
│              基础设施层           │                                  │
│           ┌───────────────┐  ┌───┴─────────────┐                   │
│           │SQLiteRepository│  │MemoryRepository │（测试专用）       │
│           │  ~/.todo.db    │  │   in-memory     │                   │
│           └───────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────────────┘
```

### 模块关系图（前端）

```
fe/
├── app/
│   ├── lib/
│   │   └── api.ts         ──── 统一 HTTP 请求封装（axios）
│   ├── components/
│   │   ├── TasksPanel.tsx  ──── 核心业务面板（列表/增/完成/删）
│   │   └── TaskFilter.tsx  ──── 状态筛选组件（全部/待办/已完成）
│   └── types/
│       └── index.ts        ──── 全局 TypeScript 类型定义
└── tests/
    ├── TaskFilter.test.tsx ──── TaskFilter 单元测试
    └── TasksPanel.test.tsx ──── TasksPanel 集成测试
```

---

## 快速上手

### 环境要求

| 工具    | 版本要求   | 安装方式              |
| ------- | ---------- | --------------------- |
| Go      | ≥ 1.22     | https://go.dev/dl/    |
| Node.js | 22.x (LTS) | https://nodejs.org/   |
| pnpm    | ≥ 9.x      | `npm install -g pnpm` |
| Git     | 任意       | https://git-scm.com/  |

> **Windows 用户**：请使用 PowerShell 或 Git Bash，不需要安装额外的 C 编译工具（本项目 SQLite 采用纯 Go 实现，无 CGO 依赖）。

---

### 后端启动

```bash
# 1. 克隆仓库
git clone https://github.com/Heuluck/NCU-SE-Group-Homework.git
cd NCU-SE-Group-Homework

# 2. 下载依赖
go mod download

# 3. 编译 CLI
go build -o gophertodo ./cmd/

# 4. 使用 CLI
./gophertodo add "完成软件工程作业"   # 添加任务
./gophertodo list                      # 查看待办列表
./gophertodo done 1                    # 完成 ID=1 的任务
./gophertodo delete 1                  # 删除 ID=1 的任务

# 5. 运行后端 API 服务器（如已实现）
go run ./cmd/server/main.go
# 默认监听 http://localhost:8080
```

**数据文件位置**：任务数据默认存储在 `~/.todo.db`（跨平台，Windows 为 `%USERPROFILE%\.todo.db`）。

---

### 前端启动

```bash
# 进入前端目录
cd fe

# 安装依赖（首次）
pnpm install

# 开发模式启动（热重载）
pnpm dev
# 访问 http://localhost:5173

# 生产构建
pnpm build
# 产物位于 fe/build/
```

---

### 使用 Mock 开发

前端可在没有真实后端的情况下，通过 Prism Mock 服务模拟 API：

```bash
cd fe

# 方式 1：使用 pnpm 脚本（推荐）
pnpm mock
# Prism Mock 监听 http://127.0.0.1:4010

# 方式 2：手动启动
npx @stoplight/prism-cli mock openapi.yaml -h 0.0.0.0 -p 4010

# 同时启动前端（另开终端）
pnpm dev
# 前端开发模式自动连接 http://127.0.0.1:4010
```

> **原理**：`fe/app/lib/api.ts` 在开发环境下默认使用 `http://127.0.0.1:4010` 作为 API 基址，Prism 读取 `openapi.yaml` 自动生成符合契约的 Mock 响应。

**环境变量配置**（可选）：

```bash
# fe/.env.local（本地覆盖，不提交 Git）
VITE_API_BASE_URL=http://localhost:8080   # 指向真实后端
```

---

## 核心模块职责说明

### 后端模块

| 模块/文件                 | 职责                                         | 关键方法                                             |
| ------------------------- | -------------------------------------------- | ---------------------------------------------------- |
| `cmd/`                    | CLI 入口，解析命令行参数                     | `cobra.Command` 子命令注册                           |
| `service/task_service.go` | **应用服务层**：编排业务逻辑，不含持久化细节 | `AddTask`, `CompleteTask`, `ListTasks`, `DeleteTask` |
| `domain/task.go`          | **领域实体**：封装业务行为，保障业务约束     | `MarkCompleted()` 状态变更                           |
| `repository/interface.go` | **仓储接口**：定义数据访问抽象（依赖倒置）   | `ITaskRepository`                                    |
| `repository/sqlite.go`    | **SQLite 实现**：生产环境数据持久化          | CRUD 操作，`~/.todo.db`                              |
| `repository/memory.go`    | **内存实现**：仅用于单元测试，无 IO 依赖     | 等同 SQLite 接口                                     |

**数据流向**：

```
用户输入 → CLI解析 → TaskAppService → ITaskRepository → SQLite/Memory
                                   ↓
                              Task实体（业务规则）
```

### 前端模块

| 模块/文件                          | 职责                                             | 说明                               |
| ---------------------------------- | ------------------------------------------------ | ---------------------------------- |
| `fe/app/lib/api.ts`                | **HTTP 客户端**：封装所有 API 请求，统一错误拦截 | 基于 axios，支持 Mock/生产环境切换 |
| `fe/app/components/TasksPanel.tsx` | **核心业务面板**：任务的增删查改完整闭环         | 含状态管理、分页、筛选、错误展示   |
| `fe/app/components/TaskFilter.tsx` | **筛选组件**：按状态筛选任务（全部/待办/已完成） | 无状态，通过 Props 传入值和回调    |
| `fe/app/types/index.ts`            | **类型定义**：Task、API 参数等 TypeScript 类型   | 全局共享，保障类型安全             |
| `fe/openapi.yaml`                  | **接口契约**：OpenAPI 3.0 规范，Prism Mock 使用  | 作为前后端联调的法定依据           |

### CI/CD 流水线

| 流水线     | 文件                            | 触发条件               | 内容                                  |
| ---------- | ------------------------------- | ---------------------- | ------------------------------------- |
| 测试流水线 | `.github/workflows/ci-test.yml` | push/PR → main         | 后端 Go 测试 + 前端 Vitest + 构建验证 |
| 部署流水线 | `.github/workflows/fe.yml`      | push → main（fe/\*\*） | 前端构建并部署到 `fe-gh-pages` 分支   |

---

## API 接口

完整接口定义见 [`fe/openapi.yaml`](./fe/openapi.yaml)，以下为快速参考：

| 方法     | 路径                   | 描述                                    |
| -------- | ---------------------- | --------------------------------------- |
| `GET`    | `/tasks`               | 获取所有任务列表                        |
| `POST`   | `/tasks`               | 创建新任务（body: `{content: string}`） |
| `GET`    | `/tasks/{id}`          | 获取指定任务详情                        |
| `POST`   | `/tasks/{id}/complete` | 标记任务为已完成                        |
| `DELETE` | `/tasks/{id}`          | 删除指定任务                            |

**Task 数据结构**：

```json
{
  "id": 1,
  "content": "完成软件工程作业",
  "status": "pending",
  "created_at": "2026-04-23T08:00:00Z",
  "completed_at": null
}
```

---

## 测试

### 后端测试（Go）

```bash
# 运行全部测试
go test ./...

# 带覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out   # 在浏览器中查看

# 运行特定测试
go test ./service/... -run TestAddTask -v
```

### 前端测试（Vitest）

```bash
cd fe

# 运行全部测试
pnpm test

# 监听模式（TDD 推荐）
pnpm test --watch

# 带覆盖率
pnpm test --coverage
```

### 当前测试覆盖情况

| 模块               | 测试文件                       | 覆盖内容                                        |
| ------------------ | ------------------------------ | ----------------------------------------------- |
| `TaskAppService`   | `service/task_service_test.go` | AddTask / CompleteTask / ListTasks / DeleteTask |
| `Task 领域实体`    | `service/task_service_test.go` | MarkCompleted 业务行为 + 幂等性                 |
| `MemoryRepository` | 集成于 service_test            | 完整 CRUD 操作                                  |
| `TaskFilter 组件`  | `fe/tests/TaskFilter.test.tsx` | 渲染 + 交互                                     |
| `TasksPanel 组件`  | `fe/tests/TasksPanel.test.tsx` | 加载 / 添加 / 完成 / 删除全流程                 |

---

## 项目结构

```
NCU-SE-Group-Homework/
│
├── .github/
│   └── workflows/
│       ├── fe.yml             # 前端部署流水线
│       └── ci-test.yml        # 测试 CI 流水线（新增）
│
├── cmd/                       # CLI 入口（cobra 命令）
│   └── main.go
│
├── domain/                    # 领域模型层
│   └── task.go                # Task 实体 + MarkCompleted()
│
├── service/                   # 应用服务层
│   ├── task_service.go        # TaskAppService
│   └── task_service_test.go   # 单元测试
│
├── repository/                # 基础设施层
│   ├── interface.go           # ITaskRepository 接口
│   ├── sqlite.go              # SQLite 实现
│   └── memory.go              # 内存实现（测试用）
│
├── fe/                        # 前端子项目
│   ├── app/
│   │   ├── components/        # React 组件
│   │   │   ├── TasksPanel.tsx
│   │   │   └── TaskFilter.tsx
│   │   ├── lib/
│   │   │   └── api.ts         # HTTP 客户端
│   │   └── types/
│   │       └── index.ts
│   ├── tests/
│   │   ├── TaskFilter.test.tsx
│   │   └── TasksPanel.test.tsx
│   ├── openapi.yaml           # OpenAPI 接口契约
│   └── package.json
│
├── specs/                     # 项目文档
│   ├── sprint-retrospective.md
│   ├── detailed-design-specification.md
│   ├── software-requirements-specification.md
│   ├── product-backlog.md
│   └── maintainability-self-assessment.md
│
├── go.mod
└── README.md
```

---

## 常见问题 FAQ

**Q: 后端编译时提示 CGO 相关错误？**  
A: 本项目使用 `modernc.org/sqlite`（纯 Go SQLite），无需 CGO。请确认 Go 版本 ≥ 1.22，并运行 `go mod tidy` 重新整理依赖。

**Q: 前端启动报 `Cannot find module` 错误？**  
A: 请先在 `fe/` 目录执行 `pnpm install`，再启动 `pnpm dev`。

**Q: Mock 服务数据每次重启都会重置？**  
A: 这是正常现象。Prism Mock 是无状态的，每次重启都从 `openapi.yaml` 中的示例数据重新生成。如需持久化数据，请启动真实后端。

**Q: 数据库文件在哪里？**  
A: 默认存储在用户主目录 `~/.todo.db`。Windows 下为 `C:\Users\<用户名>\.todo.db`。

---

_文档版本 v3.0 · 更新于 2026-04-23 · 啥队_
