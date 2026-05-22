# /help

介绍 My Harness 体系的核心架构与可用能力。

## Agent 行为

### 第一步：核心介绍

Harness 围绕四个核心支柱构建：

1. **Superpowers 体系** — 第一流程依据
  完整的编程工作流，覆盖从需求分析到代码合入的全过程（brainstorming → plan → implement → review → verify）。以严谨的流程约束和 HARD-GATE 提升开发质量，是 Harness 的流程引擎。
2. **Docs 文档体系** — 第一文档体系
  结构化的项目文档管理方案，定义文档目录布局、Git 跟踪策略、编写时机和决策矩阵。以系统化的文档产出提升项目质量和可维护性，是 Harness 的知识沉淀层。
3. **Lesson Capture 经验捕获体系** — 知识提炼引擎
  系统化捕获开发过程中反复出现的有意义问题。同一问题出现 3 次后自动提炼为独立技能并纳入 Harness 体系，将隐性经验转化为可复用知识资产。
4. **Commands 命令集合** — 统一交互入口
  覆盖开发流程、文档管理、项目管理、Git 分支等场景，所有命令遵循 Superpowers 流程优先、Docs 文档规则次之的优先级体系。

### 第二步：询问用户

输出以下引导语：

```
你想深入了解哪部分？

1. Superpowers 开发流程 — 核心流程介绍与选型理由
2. 文档体系 — 文档结构与规范说明
3. 经验捕获体系 — 问题记录、累计与技能提炼机制
4. 命令集合 — 全部命令及用途一览
```

### 第三步：根据用户选择展开

#### 1. Superpowers 开发流程

##### 核心流程

Superpowers 定义了以下关键流程节点，每个节点绑定一个 Harness 命令：


| 流程                                 | 命令           | 说明                                  |
| ---------------------------------- | ------------ | ----------------------------------- |
| **brainstorming**                  | `/design`    | 需求分析与头脑风暴，将模糊需求转化为结构化设计方案           |
| **writing-plans**                  | `/plan`      | 基于定稿设计编写实现计划，分解为可执行任务清单             |
| **executing-plans**                | `/implement` | 按计划编码实现，支持 TDD 和子代理驱动开发两种模式         |
| **requesting-code-review**         | `/review`    | 对分支改动进行系统化代码审查，发现问题回退到实现阶段          |
| **verification-before-completion** | `/accept`    | 验收前全面检查——测试通过、文档完整、流程闭环             |
| **finishing-a-development-branch** | `/accept`    | 分支完成的收尾工作：生成 release-notes、更新全局跟踪文档 |
| **systematic-debugging**           | `/debug`     | 对问题进行系统化调试，追踪根因并记录修复过程              |


**推荐主工作流**：`/design → /plan → /implement → /review → /accept`

##### 为什么选择 Superpowers？

Superpowers 是目前**唯一**同时满足以下三个条件的 Cursor 编程流程 Harness：

- **完整** — 覆盖从需求到合入的全流程，每个阶段有明确的 HARD-GATE 约束（如"未批准设计前禁止编码"），从根本上防止流程跳过
- **简练** — 流程定义精准，不引入不必要的复杂度，开发者能快速上手
- **官方面向 Cursor** — 专为 Cursor IDE 设计和维护，与 Cursor 的 Agent 模式深度适配，能稳定发挥效力

与其他方案（如 ECC）的对比：

- ECC 功能虽全面，但**官方对 Cursor 的适配明显滞后**——Codex、Claude 相关配置持续升级，而 Cursor 适配长期未更新
- ECC 的**语言/框架级别的具体编程指引实际内容很少**，侧重点在流程框架定义而非实战编码指导
- Superpowers 轻量且专注，与 Cursor 的交互模式天然契合，无适配损耗

#### 2. 文档体系

读取 `doc-project-structure` 规则和 `doc-file-definition-`* 技能后，介绍：

##### 文档体系概览

**目录布局**：

```
project-root/
├── README.md / CHANGELOG.md / version.md        ← 全局跟踪文档（Git 跟踪）
├── docs/
│   ├── ARCHITECTURE.md / API.md / DATABASE.md   ← 全局参考文档（Git 跟踪）
│   ├── superpowers/                              ← Superpowers 原生输出（不跟踪，由 Superpowers 自行管理）
│   └── <branch-id>/                              ← 分支过程文档（不跟踪，随分支清理）
│       ├── requirement.md      必写 — 需求记录
│       ├── design.md           必写 — 设计文档
│       ├── plan.md             必写 — 实现计划
│       ├── devlog.md           必写 — 开发日志（每次 commit 追加）
│       ├── accept-log.md       必写 — 计划偏离记录（修改/同意）
│       ├── review.md           选写 — 代码审查记录
│       ├── test-report.md      选写 — 测试报告
│       ├── release-notes.md    选写 — 发布说明
│       ├── issues.md           选写 — 问题跟踪
│       └── demo/               选写 — UI/流程 Demo 原型
```

**核心理念**：

- **分支 + worktree**：功能分支须在独立 git worktree 中开发（`main` 仅驻留主检出）；详见 `doc-project-structure`「强制约束 — 分支与 worktree 开发」
- **过程文档**（`docs/<branch>/`）随分支生命期维护，分支合入后即清理，不污染主分支历史
- **全局文档**（`README`、`CHANGELOG`、`ARCHITECTURE` 等）在分支合入时根据决策矩阵增量更新，持续沉淀项目知识
- **决策矩阵**明确定义了"什么业务事件触发什么文档更新"，消除维护的随意性和遗忘

**查看完整规范**：使用 `/docs-regulation-show` 查看各文档的具体格式要求，使用 `/docs-project-check` 检查当前分支的文档完整性。

#### 3. 经验捕获体系

##### 核心机制

开发过程中反复出现的有意义问题会被系统化捕获：首次出现写入问题清单，重复出现累加计数，累计 3 次后自动提炼为独立技能并纳入 Harness 体系。

##### "有意义"的判断标准

只有符合下列标准之一的问题才会被记录：


| 维度    | 说明                                |
| ----- | --------------------------------- |
| 上下文浪费 | 问题导致 Agent 大量重复理解上下文，或用户需反复描述同一诉求 |
| 生产效率  | 问题的解决能显著减少开发耗时或减少返工               |
| 代码质量  | 问题的解决能系统性避免某类 bug 或架构缺陷           |
| 部署效率  | 问题的解决能减少部署失败率或缩短部署耗时              |


##### 工作流

```
问题出现 → 首次：写入问题清单 → 重复：次数 +1
                                    │
                        次数 < 3 → 继续累计
                        次数 = 3 → 自动提炼为技能
                                    1. 创建 SKILL.md（含背景、根因、方案、场景）
                                    2. 纳入 Harness 体系
```

##### 存储结构

```
.cursor/skills/lesson-record/
├── lesson-record.md              ← 问题清单（编号/描述/次数/方案/状态）
└── experience-summary-skill/     ← 提炼后的经验技能
```

##### 触发方式

- **被动**：Agent 在会话中主动识别有意义的问题，向用户提议记录
- **显式**：`/lesson-record [问题描述]` — 结合会话上下文评估并记录问题

#### 4. 命令集合

##### 开发流程命令

遵循 Superpowers 流程，构成完整开发闭环：


| 命令             | 用途                      |
| -------------- | ----------------------- |
| `/design [主题]` | 启动完整设计流程，将需求转化为结构化设计方案  |
| `/plan`        | 基于定稿设计编写实现计划，分解为可执行任务   |
| `/implement`   | 按计划开始编码实现               |
| `/review`      | 对当前分支改动进行代码审查           |
| `/accept`      | 验收分支工作，执行分支合入合成流程       |
| `/debug [问题]`  | 对指定问题启动系统化调试            |


##### 文档管理命令


| 命令                      | 用途                      |
| ----------------------- | ----------------------- |
| `/docs-project-check`   | 检查当前分支文档完整性，逐项输出 ✅/❌    |
| `/docs-regulation-show` | 查看项目文档格式规范，可按指定文档输出完整模板 |


##### 项目管理命令


| 命令          | 用途             |
| ----------- | -------------- |
| `/run`      | 检测技术栈，启动开发服务；  |
| `/shutdown` | 识别并关闭项目相关运行中进程 |


##### Harness 验证命令


| 命令                               | 用途                                   |
| -------------------------------- | ------------------------------------ |
| `/myharness-test-project <path>` | 逻辑上切换到 `test-project/` 下的指定子项目的逻辑工作区 |


##### Git 分支命令


| 命令                           | 用途                                    |
| ---------------------------- | ------------------------------------- |
| `/git-init`                  | 初始化 Git 仓库并配置文档体系（.gitignore、初始提交）    |
| `/git-load [远端]`             | 从远端同步项目到本地（仅空项目可用）                    |
| `/git-current-branch-info`   | 查看当前分支完整信息（含 worktree 绑定状态、文档、阶段推断）      |
| `/git-current-branch-commit` | 提交变更，自动更新 devlog 并生成规范 commit message |
| `/git-current-branch-push`   | 推送当前分支到远端                             |
| `/git-current-branch-merge`  | 合并当前分支到主分支（前置检查 accept 流程）            |
| `/git-current-branch-delete` | 删除当前分支并移除绑定 worktree（main 禁止删除，二次确认）     |
| `/git-list-branch-all`       | 列出所有未 accept 的文档管理分支                  |
| `/git-list-branch-bug`       | 列出 Bug 类型分支，标注未解决问题数                  |
| `/git-list-branch-rp`        | 列出需求类型分支，标注当前所处阶段                     |
| `/help`                      | 显示本帮助信息                               |


