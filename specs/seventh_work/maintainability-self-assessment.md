# 可维护性五因素自评报告

**项目名称**: GopherTodo — 全栈待办事项系统  
**团队**: 啥队（南昌大学计科软件工程小组）  
**成员**: 陆永祥 · 邓枭 · 马骏  
**评估版本**: Sprint 3 当前代码快照  
**评估日期**: 2026-04-23  
**评估依据**: 教材 §8.4.1 软件可维护性五因素

---

## 一、评估说明

### 1.1 评分标准

| 分数 | 等级 | 含义 |
|------|------|------|
| 5 | 优秀 | 完全满足标准，无明显缺陷 |
| 4 | 良好 | 基本满足标准，少量改进空间 |
| 3 | 中等 | 部分满足标准，有明显待改进点 |
| 2 | 较差 | 明显不足，需要系统性改进 |
| 1 | 很差 | 严重缺失，需立即整治 |

### 1.2 评估范围

本次评估覆盖两个子系统：

- **后端（GopherTodo CLI）**：Golang + SQLite + DDD分层架构
- **前端（GopherTodo FE）**：React + TypeScript + Ant Design + Vite

---

## 二、五因素评分总览

| 因素 | 后端得分 | 前端得分 | 综合得分 | 优先级 |
|------|---------|---------|---------|--------|
| 可理解性（Understandability） | 3 | 4 | **3.5** | 🟡 中 |
| 可测试性（Testability） | 2 | 3 | **2.5** | 🔴 低 |
| 可修改性（Modifiability） | 4 | 3 | **3.5** | 🟡 中 |
| 可移植性（Portability） | 3 | 4 | **3.5** | 🟡 中 |
| 可重用性（Reusability） | 3 | 3 | **3.0** | 🟡 中 |

**总体加权得分**: **3.2 / 5**（需重点提升可测试性）

---

## 三、各因素详细评估

### 3.1 可理解性（Understandability）— 综合 3.5 分

#### 后端评分：3 / 5

**评估依据**：

| 维度 | 现状 | 问题 |
|------|------|------|
| 命名规范 | Go 命名风格基本规范 | 部分内部函数缺少注释 |
| 架构清晰度 | DDD分层明确：CLI→AppService→Domain→Repository | specs文档与代码存在一定滞后 |
| 注释覆盖率 | 公共方法缺少GoDoc格式注释 | 约50%的exported函数无注释 |
| 错误信息 | 错误类型初步定义，但信息不够详细 | ErrTaskNotFound等缺少上下文 |

**代码示例（现状）**：
```go
// 现有代码缺少注释
func (s *TaskAppService) CompleteTask(id int) error {
    task, err := s.repo.FindById(id)
    if err != nil {
        return err
    }
    // ...
}
```

#### 前端评分：4 / 5

**评估依据**：

| 维度 | 现状 | 评价 |
|------|------|------|
| 组件结构 | TasksPanel.tsx 功能清晰，职责明确 | ✅ 良好 |
| 类型定义 | TypeScript类型完备，接口定义清晰 | ✅ 良好 |
| API封装 | api.ts模块设计简洁，函数名见名知意 | ✅ 良好 |
| OpenAPI文档 | openapi.yaml覆盖所有接口，有示例数据 | ✅ 良好 |
| 注释 | 业务逻辑注释偶有缺失 | ⚠️ 一般 |

---

### 3.2 可测试性（Testability）— 综合 2.5 分 🔴

#### 后端评分：2 / 5

**评估依据**：

| 维度 | 现状 | 问题 |
|------|------|------|
| 单元测试覆盖 | 当前**几乎无后端单元测试** | 严重不足 |
| Mock能力 | ITaskRepository接口已定义 | MemoryRepository未落地 |
| 可测函数比例 | AppService层方法可测试 | 未编写测试文件 |
| CI集成 | GitHub Actions仅跑前端构建 | 无后端测试流水线 |

#### 前端评分：3 / 5

**评估依据**：

| 维度 | 现状 | 问题 |
|------|------|------|
| 组件测试 | TaskFilter.test.tsx 覆盖基础渲染 | 仅1个测试套件 |
| API测试 | api.ts 无测试 | Mock接口无覆盖 |
| 业务逻辑测试 | TasksPanel交互逻辑无测试 | 缺乏集成测试 |
| 测试工具链 | React Testing Library 已引入 | 测试框架就绪但用量少 |

---

### 3.3 可修改性（Modifiability）— 综合 3.5 分

#### 后端评分：4 / 5

**评估依据**：

| 维度 | 现状 | 评价 |
|------|------|------|
| 模块耦合 | Sprint 2已完成DDD分层解耦 | ✅ 良好 |
| 接口设计 | ITaskRepository接口隔离存储层 | ✅ 良好 |
| 修改影响范围 | 存储实现变更只需替换Repository | ✅ 良好 |
| 技术债务 | 少量硬编码配置（DB路径） | ⚠️ 轻微 |

#### 前端评分：3 / 5

**评估依据**：

| 维度 | 现状 | 问题 |
|------|------|------|
| 组件粒度 | TasksPanel功能相对集中 | 可进一步拆分子组件 |
| 状态管理 | 使用React useState，无全局状态 | 随业务增长可能需要引入Context/Zustand |
| API层解耦 | api.ts独立封装，修改集中 | ✅ 良好 |
| 样式可维护 | Tailwind CSS直接书写，局部耦合 | ⚠️ 一般 |

---

### 3.4 可移植性（Portability）— 综合 3.5 分

#### 后端评分：3 / 5

**评估依据**：

| 维度 | 现状 | 问题 |
|------|------|------|
| 跨平台编译 | Go单二进制，理论支持跨平台 | ⚠️ DB路径依赖os.UserHomeDir()未完全验证 |
| CGO依赖 | 使用modernc.org/sqlite（无CGO） | ✅ 良好 |
| 环境配置 | 数据库路径有硬编码风险 | 配置模块尚未抽象 |
| 容器化 | 暂无Dockerfile | 可以优化 |

#### 前端评分：4 / 5

**评估依据**：

| 维度 | 现状 | 评价 |
|------|------|------|
| 浏览器兼容 | Vite + React标准构建，兼容现代浏览器 | ✅ 良好 |
| API地址配置 | 通过VITE_API_BASE_URL环境变量控制 | ✅ 良好 |
| 部署方式 | GitHub Actions自动部署到gh-pages | ✅ 良好 |
| 移动端适配 | Ant Design响应式，但未专门优化 | ⚠️ 一般 |

---

### 3.5 可重用性（Reusability）— 综合 3.0 分

#### 后端评分：3 / 5

**评估依据**：

| 维度 | 现状 | 问题 |
|------|------|------|
| 领域模型 | Task实体设计良好，可复用于REST API | ✅ 良好 |
| 仓储接口 | ITaskRepository可重用于不同存储实现 | ✅ 良好 |
| 应用服务 | TaskAppService耦合程度适中 | ⚠️ 可进一步抽象 |
| 包结构 | 包结构设计尚未严格模块化 | 改进空间较大 |

#### 前端评分：3 / 5

**评估依据**：

| 维度 | 现状 | 评价 |
|------|------|------|
| 通用组件 | TaskFilter可独立复用 | ✅ 良好 |
| API模块 | api.ts作为独立模块可复用 | ✅ 良好 |
| 业务组件 | TasksPanel业务特定性较强 | ⚠️ 复用性有限 |
| 工具函数 | 缺少通用工具函数库 | ⚠️ 有改进空间 |

---

## 四、低于 3 分因素的改进方案

### 4.1 可测试性提升方案（最高优先级 P0）

#### 问题根因

后端单元测试完全缺失，前端测试覆盖度极低，CI流水线中没有自动化测试环节。

#### 改进方案

**方案一：后端单元测试体系建设**

```go
// 步骤1: 实现MemoryRepository（用于单元测试mock）
type MemoryRepository struct {
    tasks  map[int]*Task
    nextID int
    mu     sync.RWMutex
}

func (r *MemoryRepository) Save(task *Task) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.nextID++
    task.ID = r.nextID
    r.tasks[task.ID] = task
    return nil
}
// 其余方法实现...

// 步骤2: 编写AppService单元测试
func TestAddTask(t *testing.T) {
    repo := NewMemoryRepository()
    svc  := NewTaskAppService(repo)

    task, err := svc.AddTask("测试任务")
    assert.NoError(t, err)
    assert.Equal(t, "测试任务", task.Content)
    assert.Equal(t, Pending, task.Status)
}

func TestCompleteTask(t *testing.T) {
    repo := NewMemoryRepository()
    svc  := NewTaskAppService(repo)
    task, _ := svc.AddTask("任务A")

    err := svc.CompleteTask(task.ID)
    assert.NoError(t, err)

    updated, _ := svc.GetTask(task.ID)
    assert.Equal(t, Completed, updated.Status)
    assert.NotNil(t, updated.CompletedAt)
}

func TestCompleteTask_NotFound(t *testing.T) {
    repo := NewMemoryRepository()
    svc  := NewTaskAppService(repo)

    err := svc.CompleteTask(999)
    assert.ErrorIs(t, err, ErrTaskNotFound)
}
```

**方案二：前端测试补全**

```tsx
// TasksPanel交互测试
describe('TasksPanel - 交互逻辑', () => {
  it('应能成功添加任务', async () => {
    render(<TasksPanel />);
    const input = screen.getByPlaceholderText('输入新任务...');
    await userEvent.type(input, '测试任务{enter}');
    await waitFor(() => {
      expect(screen.getByText('测试任务')).toBeInTheDocument();
    });
  });
});
```

**方案三：CI流水线集成测试**

```yaml
# .github/workflows/test.yml
name: Run Tests
on: [push, pull_request]
jobs:
  backend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: go test ./... -v -coverprofile=coverage.out
      - run: go tool cover -func=coverage.out
  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '22' }
      - run: cd fe && npm i -g pnpm && pnpm install && pnpm test
```

**执行计划**：

| 步骤 | 负责人 | 工期 |
|------|--------|------|
| 实现MemoryRepository | 马骏 | 1天 |
| 编写后端3个核心测试 | 马骏 | 1天 |
| 补全前端TasksPanel测试 | 陆永祥 | 1天 |
| 配置test CI流水线 | 邓枭 | 0.5天 |

---

## 五、可维护性改进路线图

```
Sprint 3 (当前)       Sprint 4              Sprint 5
─────────────         ─────────────         ─────────────
✅ DDD分层重构         ✅ 单元测试 ≥80%      配置管理模块
                      ✅ CI测试流水线        日志可观测性
                      GoDoc注释补全         TUI/移动端适配
                      MemoryRepository      通用组件库
```

---

**评估人**: 啥队全员  
**下次评估**: Sprint 4 结束后
