# Harness 组件依赖关系图

Harness 各组件之间存在联动关系——修改 A 必须同步更新 B。本文档是 Harness 的**依赖关系唯一信源**，Agent 在执行任何修改操作前应先查阅本文档，修改完成后按本文档完成所有联动更新。

依赖分为两类：
- **泛性依赖**：按组件类型定义，凡属于该类型的组件发生变更，均需执行对应联动
- **指定依赖**：特定组件之间存在硬性引用关系（A 引用了 B 的名称或行为），变更 B 时必须同步更新 A

---

## 一、泛性依赖（按组件类型）

凡新增、删除、重命名或修改行为描述，均须同步更新 `README.md`。README.md 是 Harness 的聚合展示文档，所有组件变更都映射到此。

| 组件类型 | 变更操作 | README.md 中的更新位置 | 具体动作 |
|---------|---------|---------------------|---------|
| **Rule** | 新增 | `## 架构总览` 树形图 | rules/ 下新增一行，计数 +1 |
| Rule | 删除 | `## 架构总览` 树形图 | 移除对应行，计数 -1 |
| Rule | 重命名 | `## 架构总览` 树形图 | 更新文件名和描述 |
| Rule | 修改行为 | `## 架构总览` 树形图 | 更新描述文字（如行为变化导致描述不准确） |
| **Skill** | 新增 | `## 架构总览` 树形图 | skills/ 下新增一行，计数 +1 |
| Skill | 删除 | `## 架构总览` 树形图 | 移除对应行，计数 -1 |
| Skill | 重命名 | `## 架构总览` 树形图 | 更新目录名 |
| Skill（子技能） | 新增/删除 | `## 架构总览` 树形图 + 父技能 `SKILL.md` | 父技能目录下更新子技能行；父技能「子技能目录」表同步 |
| **Command** | 新增 | `## 架构总览` 树形图 + `## Commands` 表格 | 树形图新增一行，计数 +1；对应分类表格新增一行 |
| Command | 删除 | `## 架构总览` 树形图 + `## Commands` 表格 | 树形图移除对应行，计数 -1；对应分类表格移除对应行 |
| Command | 重命名 | `## 架构总览` 树形图 + `## Commands` 表格 | 两处同步更新名称 |
| Command | 修改行为 | `## Commands` 表格 | 更新说明文字 |
| **Subagent** | 新增 | `## 架构总览` 树形图 | subagents/ 下新增内容 |
| Subagent | 删除 | `## 架构总览` 树形图 | 移除对应内容 |
| **Hook** | 新增 | `## 架构总览` 树形图 | hooks.json 和 hooks/ 新增内容 |
| Hook | 删除 | `## 架构总览` 树形图 | 移除对应内容 |

> **重命名时的额外操作**：除更新 README.md 外，还必须 Grep 搜索旧名称（范围为 `README.md` 和 `.cursor/` 目录），将全部引用处替换为新名称。

---

## 二、指定依赖（按具体组件）

以下列出存在硬性引用关系的组件对。关系方向为：**左侧组件被右侧组件引用**。当左侧组件发生变更（重命名、改行为、改结构）时，必须检查并同步更新右侧列出的所有组件。

### Rules → 引用者

| 被引用的 Rule | 必须同步更新的引用者 | 引用方式 |
|-------------|-------------------|---------|
| `doc-project-structure` | 全部 15 个 `doc-file-definition-*` 技能 | 技能引用其定义的目录布局、文档触发时机和决策矩阵 |
| | `/design`、`/plan`、`/implement`、`/review`、`/accept`、`/refuse`、`/debug` 命令 | 命令引用其定义的分支文档结构 |
| | `/docs-project-check`、`/docs-regulation-show` 命令 | 命令以其定义的文档清单、必写/选写、创建/触发时机与决策矩阵为核心；正文格式由 `doc-file-definition-*` 技能提供 |
| | `/git-init`、`/git-current-branch-commit` 命令 | 命令引用其定义的文档跟踪策略和 devlog 写入时机 |
| | `architecture-design` 技能 | 技能产出写入 `docs/<branch>/design.md`，并按决策矩阵联动 ARCHITECTURE / API / DATABASE 等全局文档 |
| | `my-harness-build` 规则 | 信源唯一性表中将其列为文档架构唯一源头 |
| | `priority-resolution` 规则 | 第二优先级裁决引用该规则 |
| `priority-resolution` | `my-harness-build` 规则 | 优先级裁决原则引用 |
| | `architecture-design` 技能 | 技能在与 Superpowers / `doc-project-structure` 冲突时按其裁决顺序 |
| `brainstorming-visual-enhancement` | `/design` 命令 | `/design` 流程中插入的可视化触发条件 |
| | `doc-file-definition-design` 技能 | 场景 UI 时前端设计 AskQuestion 与 `design.md`「前端设计」节格式 |
| | `brainstorming-architecture-thinking-enhancement` 规则 | 触发条件叠加时两规则并列生效（视觉 / 架构两维度） |
| `brainstorming-architecture-thinking-enhancement` | `/design` 命令 | `/design` 流程中插入的架构思考触发条件 |
| | `architecture-design` 技能 | 规则作为插件接入 Superpowers brainstorming，强制加载本技能 |
| | `my-harness-build` 规则 | 信源唯一性表将其列为"架构思考接入方式"唯一源头 |
| `lesson-capture` | `.cursor/skills/lesson-record/lesson-record.md` | 问题清单表格结构、提炼阈值（3 次）均由该规则定义 |
| `concise-communication` | `my-harness-build` 规则 | 信源唯一性表将其列为沟通风格唯一源头 |

### Skills → 引用者

| 被引用的 Skill | 必须同步更新的引用者 | 引用方式 |
|-------------|-------------------|---------|
| `frontend-design` | `brainstorming-visual-enhancement` 规则 | UI 场景按 `design.md`「前端设计」读取该技能并路由 style/frame |
| | `my-harness-build` 规则 | 信源唯一性表将其列为前端设计唯一源头 |
| | `README.md` | 架构总览中 frontend-design 目录结构、参考索引与 ui-frame 子技能 |
| `preview-html` | `brainstorming-visual-enhancement` 规则 | Demo/线框 HTML 预览须读取并执行该技能 |
| | `my-harness-build` 规则 | 信源唯一性表将其列为 HTML 预览唯一源头 |
| | `README.md` | 架构总览中 preview-html 目录 |
| `frontend-design` 子技能或参考索引变更 | `frontend-design/SKILL.md` | 约束式子技能表或「参考索引」需同步 |
| | `README.md` | 架构总览中子技能行需同步 |
| `doc-file-definition-design` | `architecture-design` 技能 | 整体格式以本技能为准；`architecture-design` 仅填充其「架构设计」节内部，不重定义外层结构 |
| `doc-file-definition-architecture` | `architecture-design` 技能 | 全局 ARCHITECTURE.md 的「章节增量来源映射」以 `architecture-design` 与 Superpowers spec 为增量来源；映射变更需同步 |
| | `/accept` 命令（隐含） | 合入时按本技能映射表回流分支级增量 |
| `doc-file-definition-api` | `architecture-design` 技能 | API.md「通用约定」的错误模型 / 幂等 / 版本策略必须与 `architecture-design` 接口契约一致；接口契约变更需同步 |
| `architecture-design` | `brainstorming-architecture-thinking-enhancement` 规则 | 规则作为插件接入 Superpowers brainstorming，强制加载该技能并按其工作流产出 |
| | `my-harness-build` 规则 | 信源唯一性表将其列为技术架构设计唯一源头 |
| | `README.md` | 架构总览中 architecture-design 目录结构与引用文件 |
| `architecture-design` 内部文件 (`workflow-greenfield.md` / `workflow-iteration.md` / `reference.md`) | `architecture-design/SKILL.md` | 入口对引用文件的路由表必须同步；新增 / 重命名 / 删除引用文件需同步更新入口表 |
| | `brainstorming-architecture-thinking-enhancement` 规则 | 规则第一节"缺口表"与第四节"产物表"按文件名指引（`workflow-greenfield.md` / `workflow-iteration.md` / `reference.md`），文件改名需同步 |

### Commands → 引用者

| 被引用的 Command | 必须同步更新的引用者 | 引用方式 |
|-----------------|-------------------|---------|
| `/plan` | `/design` 命令 | `/design` 完成后引导下一步为 `/plan` |
| `/implement` | `/plan` 命令 | `/plan` 完成后引导下一步为 `/implement` |
| `/review` | `/implement` 命令 | `/implement` 完成后引导下一步为 `/review` |
| `/accept` | `/review` 命令 | `/review` 通过后引导下一步为 `/accept` |
| `/implement` | `/review` 命令 | `/review` 发现问题后回退到 `/implement` |
| `/docs-project-check` | `/accept` 命令 | `/accept` 流程中调用 `/docs-project-check` |
| `/accept` | `/git-current-branch-merge` 命令 | merge 前置检查 `/accept` 是否完成 |
| `my-harness-build` 规则（test-project 小节） | `/myharness-test-project` 命令 | 规则要求验证时通过该命令激活逻辑工作区 |
| `/myharness-test-project` | `my-harness-build` 规则 | 命令仅引用工作区约束与跨项目隔离，不复述规则正文 |

---

## 三、联动更新执行顺序

当一次修改涉及多项联动时，按以下顺序执行以降低遗漏风险：

1. **先完成主体变更**：添加/修改/删除目标文件本身
2. **更新 README.md**：架构总览（树形图、计数）→ 命令表格
3. **处理交叉引用**：Grep 搜索旧名称 → 逐文件更新引用处
4. **检查指定依赖**：对照"二、指定依赖"表格，逐项检查并更新引用者

---

## 四、示例

### 示例 1：新增一个 Command

```
1. 创建 .cursor/commands/foo.md
2. README.md：
   - 架构总览树形图中 commands/ 下新增一行，计数 +1
   - "## Commands" 对应分类表格中新增一行
3. 如该 Command 引用了已有 Command（如引导下一步），检查"二、指定依赖"是否需要登记新的引用关系
```

### 示例 2：重命名 `doc-project-structure` 规则

```
1. 重命名 .cursor/rules/doc-project-structure.mdc → new-name.mdc
2. README.md：更新架构总览树形图中的文件名和描述
3. Grep "doc-project-structure" 在 README.md 和 .cursor/ 下的所有出现 → 共约 30 处
4. 逐文件将 "doc-project-structure" 替换为 "new-name"
5. 对照"二、指定依赖"中 Rules 表格，确认全部引用者已更新
```

### 示例 3：删除一个 Skill

```
1. 删除 .cursor/skills/foo/SKILL.md 及目录
2. README.md：架构总览树形图中移除对应行，计数 -1
3. Grep 搜索该 Skill 名称，清理所有残留引用
```
