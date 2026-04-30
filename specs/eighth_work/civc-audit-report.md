# CIVC 自评审计报告

## 一、CIVC 四阀门框架概述

CIVC（Constraints-Informed-Verified-Corrected）是 AI 辅助编程的安全框架，通过四个阀门控制 AI 生成代码的质量和安全性：

| 阀门 | 名称 | 作用 |
|------|------|------|
| **C** | Constraints（约束） | 定义 AI 能访问/修改的范围 |
| **I** | Informed（告知） | 提供足够的上下文让 AI 理解代码 |
| **V** | Verified（验证） | 通过自动化测试验证 AI 生成的代码 |
| **C** | Corrected（纠正） | 出错时能快速回滚到安全状态 |

---

## 二、当前状态审计

### (a) 约束（Constraints）— ⚠️ 薄弱

#### 当前状态
- **访问范围**：AI 可访问项目中所有文件（Go/TypeScript/配置文件）
- **修改权限**：无沙盒隔离，AI 可修改任何文件
- **敏感区域**：无明确保护机制，可能误修改核心 domain 逻辑

#### 存在问题
```
1. 无 .claudeignore 或 .cursorignore 等文件
2. domain/task.go 是核心领域逻辑，但 AI 可直接修改
3. 缺乏文件级别的修改权限控制
4. 未配置 API 访问限制（如不允许调用外部 API）
```

#### 改进方案
```markdown
# 1. 创建 .claudeignore 文件
# 隔离敏感文件和目录

# 保护核心 domain 逻辑
internal/domain/task.go
internal/domain/task_test.go

# 保护 CI 配置（防止 AI 绕过测试）
.github/workflows/*

# 保护文档
specs/
README*.md

# 2. 在 AGENTS.md 中添加约束声明
## 禁止操作清单
- 禁止修改 domain/task.go 的业务规则
- 禁止修改 .github/workflows/ 下的 CI 配置
```

---

### (b) 告知（Informed）— ⚠️ 需改进

#### 当前状态
- ✅ 有 OpenAPI 文档（`fe/openapi.yaml`）
- ✅ 有 README 项目概述
- ❌ domain 层缺乏业务规则注释
- ❌ Service 层方法缺少文档注释

#### 存在问题
```
1. NewTask 为何要 TrimSpace？缺乏注释说明
2. MarkCompleted 的幂等性设计未说明
3. 错误 sentinel 的使用场景未解释
4. AI 难以理解"任务完成"的状态机转换
```

#### 改进方案
```go
// domain/task.go

// NewTask creates a new task with the given content.
// Content is trimmed of leading/trailing whitespace.
// Returns ErrEmptyContent if the trimmed content is empty.
// The task is created in StatusPending state.
func NewTask(content string, now time.Time) (*Task, error) {
    content = strings.TrimSpace(content)
    if content == "" {
        return nil, ErrEmptyContent
    }
    return &Task{...}, nil
}

// MarkCompleted marks the task as completed.
// This operation is idempotent: calling MarkCompleted on an already
// completed task returns ErrAlreadyDone but does not modify the task state.
// The CompletedAt timestamp is set to the current time.
func (t *Task) MarkCompleted(now time.Time) error {
    if t.Status == StatusCompleted {
        return ErrAlreadyDone
    }
    // ... implementation
}
```

---

### (c) 验证（Verified）— ✅ 良好

#### 当前状态
- ✅ CI 流水线配置完整（`.github/workflows/ci-test.yml`）
- ✅ Go 测试 + 覆盖率检查（阈值 60%）
- ✅ 前端 Vitest 测试
- ⚠️ 覆盖率阈值 60% 低于目标 80%

#### 存在问题
```
1. CI 覆盖率阈值 60% 低于作业要求的 80%
2. domain/service/repository 分层覆盖率未单独检查
3. PR 前置检查可能缺失（需确认是否强制要求）
```

#### 改进方案
```yaml
# .github/workflows/ci-test.yml

- name: Check coverage threshold (≥ 80%)
  run: |
    # 分层检查覆盖率
    for pkg in "domain" "service" "repository"; do
      COVERAGE=$(go tool cover -func=coverage.out | grep "$pkg" | awk '{print $3}' | sed 's/%//')
      echo "$pkg coverage: ${COVERAGE}%"
      if (( $(echo "$COVERAGE < 80" | bc -l) )); then
        echo "::error::$pkg coverage ${COVERAGE}% is below threshold 80%"
        exit 1
      fi
    done
```

---

### (d) 纠正（Corrected）— ⚠️ 需改进

#### 当前状态
- ✅ Git 版本控制，可手动回滚
- ❌ 无一键回滚 workflow
- ❌ 缺乏 revert 脚本或自动化回滚机制

#### 存在问题
```
1. AI 生成错误代码后，需要手动 git revert
2. 无自动化的 rollback workflow
3. 部署后的回滚需要手动干预
```

#### 改进方案
```yaml
# .github/workflows/rollback.yml
name: Rollback

on:
  workflow_dispatch:
    inputs:
      version:
        description: 'Version/Tag to rollback to'
        required: true

jobs:
  rollback:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Rollback to specified version
        run: |
          git reset --hard ${{ github.event.inputs.version }}
          git push --force origin main

      - name: Notify team
        uses: slackapi/slack-github-action@v1
        with:
          channel-id: 'C0XXXXXX'
          payload: |
            {
              "text": "🔄 Rollback triggered to ${{ github.event.inputs.version }}"
            }
```

---

## 三、改进优先级汇总

| 阀门 | 当前状态 | 优先级 | 改进方案 | 预计工时 |
|------|---------|--------|---------|---------|
| **(a) 约束** | ⚠️ 薄弱 | P1 | 添加 .claudeignore + AGENTS.md 约束 | 0.5h |
| **(b) 告知** | ⚠️ 需改进 | P2 | 添加 domain/service 层注释 | 1h |
| **(c) 验证** | ✅ 良好 | P1 | 提升 CI 覆盖率阈值至 80% | 0.5h |
| **(d) 纠正** | ⚠️ 需改进 | P2 | 添加 rollback workflow | 1h |

---

## 四、改进实施计划

### 第一阶段（立即实施）
1. 创建 `.claudeignore` 文件，保护核心文件
2. 更新 `AGENTS.md`，添加禁止操作清单
3. 提升 CI 覆盖率阈值至 80%

### 第二阶段（本周内完成）
4. 为 domain/task.go 添加业务规则注释
5. 为 service/task_service.go 添加方法注释
6. 编写 rollback workflow

### 第三阶段（后续迭代）
7. 添加 pre-commit hook 自动化检查
8. 引入 AI 代码审查工具（如 CodeRabbit）

---

## 五、验证方法

改进后，使用以下方法验证 CIVC 框架有效性：

```bash
# 1. 验证约束
# - 创建一个新的 AI 会话，仅提供 AGENTS.md
# - 询问 AI "修改 domain/task.go 的 NewTask 方法"
# - 预期：AI 应拒绝或明确表示需要确认

# 2. 验证告知
# - 检查 domain/task.go 是否有完整的文档注释
# - 检查 service/task_service.go 是否有方法说明

# 3. 验证测试
# - 运行 go test -coverprofile=coverage.out ./...
# - 检查 domain/service/repository 覆盖率是否 ≥80%

# 4. 验证纠正
# - 模拟一次错误提交
# - 触发 rollback workflow
# - 验证代码是否回滚到指定版本
```

---

## 六、结论

通过 CIVC 四阀门审计，我们发现：

1. **约束（Constraints）**：最薄弱，需要立即添加文件保护机制
2. **告知（Informed）**：基本可用，但缺少内联注释
3. **验证（Verified）**：CI 配置良好，但覆盖率阈值偏低
4. **纠正（Corrected）**：依赖 Git 手动回滚，缺乏自动化

建议按优先级实施改进，确保 AI 辅助编程的安全性和代码质量。
