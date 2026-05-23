# Cursor Harness — Superpowers 增强体系

## 项目定位

Harness 是对 **Superpowers 插件体系** 的增强层，通过在项目级别叠加 Rules、Skills、Commands、Subagents、Hooks，弥补 Superpowers 在文档管理、触发方式、可视化沟通三个方面的不足。

**三体系分工**（Superpowers 开发流程 / Doc 文档架构 / Commands 触发器）见 `my-harness-build` 规则 §1「三体系架构」。

**核心约束：不改动 Superpowers 任何内容，仅在本项目范围内新增/修改/删除。**

---

## 架构总览

```
my_harness/
├── README.md
├── .gitignore
└── .cursor/
    ├── rules/                          8 个规则
    │   ├── my-harness-build.mdc         Harness 设计原理、操作边界与扩展约束
    │   ├── priority-resolution.mdc     优先级裁决
    │   ├── concise-communication.mdc   精简沟通
    │   ├── doc-project-structure.mdc   软件开发项目的文档架构
    │   ├── brainstorming-visual-enhancement.mdc  可视化增强
    │   ├── brainstorming-architecture-thinking-enhancement.mdc  架构思考增强（接入 architecture-design）
    │   ├── using-subagent-enhancement.mdc  Subagent 使用增强（默认 Subagent-Driven 派发）
    │   └── lesson-capture.mdc          经验捕获与技能提炼
    ├── skills/                         20 个技能目录
    │   ├── doc-file-definition-requirement/
    │   ├── doc-file-definition-design/   # design.md 含「前端设计」节（UI 必填）
    │   ├── doc-file-definition-plan/
    │   ├── doc-file-definition-devlog/
    │   ├── doc-file-definition-accept-log/
    │   ├── doc-file-definition-review/
    │   ├── doc-file-definition-test-report/
    │   ├── doc-file-definition-release-notes/
    │   ├── doc-file-definition-issues/
    │   ├── doc-file-definition-readme/
    │   ├── doc-file-definition-changelog/
    │   ├── doc-file-definition-version/
    │   ├── doc-file-definition-architecture/
    │   ├── doc-file-definition-api/
    │   ├── doc-file-definition-database/
    │   ├── lesson-record/              经验记录（问题清单 + 提炼技能）
    │   │   ├── lesson-record.md        问题清单表格
    │   │   └── experience-summary-skill/  提炼后的经验技能
    │   ├── architecture-design/         架构设计指导（0-1 初建 / 1-N 迭代，含调整判断与迁移策略）
    │   │   ├── SKILL.md                 入口：场景路由 + 核心原则 + AI 行为约束
    │   │   ├── workflow-greenfield.md   0-1 初建完整工作流
    │   │   ├── workflow-iteration.md    1-N 迭代工作流（含 14 项调整判断）
    │   │   └── reference.md             架构要素参考 + 设计文档模板 + 检查清单
    │   ├── tech-stack-guidance/         组织技术栈与系统开发模式（工具级/平台级）
    │   ├── frontend-design/             前端 style+frame 路由（单技能）
    │   │   ├── SKILL.md                 模式路由、参考索引、子技能表
    │   │   ├── free-style/              开放式美学子技能（无 ui-frame 库时）
    │   │   ├── ui-style/                参考数据（<slug>/DESIGN.md）
    │   │   └── ui-frame/                约束式子技能（antd 等）
    │   └── preview-html/                本地 HTML HTTP 预览（daemon、多 root、Cursor 内置浏览器）
    │       ├── SKILL.md
    │       └── scripts/demo-preview-server.cjs
    └── commands/                       24 个命令
        ├── design.md                   /design
        ├── plan.md                     /plan
        ├── implement.md                /implement
        ├── review.md                   /review
        ├── debug.md                    /debug
        ├── accept.md                   /accept
        ├── lesson-record.md            /lesson-record
        ├── docs-project-check.md       /docs-project-check
        ├── docs-regulation-show.md     /docs-regulation-show
        ├── run.md                      /run
        ├── shutdown.md                 /shutdown
        ├── git-list-branch-all.md      /git-list-branch-all
        ├── git-list-branch-bug.md      /git-list-branch-bug
        ├── git-list-branch-rp.md       /git-list-branch-rp
        ├── git-current-branch-info.md  /git-current-branch-info
        ├── git-current-branch-commit.md /git-current-branch-commit
        ├── git-current-branch-merge.md /git-current-branch-merge
        ├── git-current-branch-push.md  /git-current-branch-push
        ├── git-init.md                  /git-init
        ├── git-load.md                  /git-load
        ├── git-current-branch-delete.md   /git-current-branch-delete
        ├── help.md                     /help
        ├── gogogo.md                    /gogogo
        ├── myharness-test-project.md   /myharness-test-project
```

---

## 文档体系

由 `doc-project-structure` 规则定义。分支过程文档 9 个（5 必写 + 4 选写），全局跟踪文档 6 个；功能分支须绑定独立 git worktree（`main` 仅驻留主检出）。决策矩阵映射业务事件到文档动作。详见该规则。

---

## 经验捕获体系

由 `lesson-capture` 规则定义。开发过程中反复出现的有意义问题将被系统化捕获：首次出现写入问题清单，重复出现累加计数，累计 3 次后自动提炼为独立技能并纳入 Harness 体系。

### "有意义"的判断标准

只有以下四类问题会被记录：

| 维度 | 说明 |
|------|------|
| 上下文浪费 | 问题导致 Agent 大量重复理解上下文，或用户需反复描述同一诉求 |
| 生产效率 | 问题的解决能显著减少开发耗时或减少返工 |
| 代码质量 | 问题的解决能系统性避免某类 bug 或架构缺陷 |
| 部署效率 | 问题的解决能减少部署失败率或缩短部署耗时 |

### 工作流

```
问题出现（符合四类标准之一）
  │
  ├─ 首次 → 写入问题清单（编号 Lxxx，记录问题+方案+时间，次数=1）
  │
  ├─ 重复 → 次数 +1
  │         │
  │         ├─ 次数 < 3 → 仅更新次数
  │         │
  │         └─ 次数 = 3 → 自动提炼：
  │              1. 创建独立技能（experience-summary-skill/<skill-name>/SKILL.md）
  │              2. 解决方式迁入技能文件
  │              3. 问题清单清空解决方式、标记已提炼、记录技能名
  │              4. 纳入 Harness 体系架构
  │
  └─ 不具意义 → 不记录
```

### 存储结构

```
.cursor/skills/lesson-record/
├── lesson-record.md              ← 问题清单表格（7 列，含编号/描述/次数/方案/提炼状态）
└── experience-summary-skill/     ← 提炼后的经验技能
    └── <skill-name>/
        └── SKILL.md
```

---

## Commands

### 开发流程
| 命令 | 说明 |
|------|------|
| `/design [主题]` | 启动设计流程 |
| `/plan` | 编写实现计划 |
| `/implement` | 开始编码 |
| `/review` | 审查代码 |
| `/debug [问题]` | 系统化调试 |
| `/accept` | 验收合入 |

### 经验捕获
| 命令 | 说明 |
|------|------|
| `/lesson-record [问题]` | 记录有意义问题到经验捕获清单，累计达阈值自动提炼为技能 |

### 文档管理
| 命令 | 说明 |
|------|------|
| `/docs-project-check` | 检查分支文档完整性 |
| `/docs-regulation-show` | 查看文档格式规范 |

### 项目管理
| 命令 | 说明 |
|------|------|
| `/run` | 启动项目开发服务 |
| `/shutdown` | 关闭项目相关进程 |

### Harness 验证
| 命令 | 说明 |
|------|------|
| `/myharness-test-project <path>` | 激活 `test-project/` 下子项目的逻辑工作区，开始针对该项目的 Harness 验证开发（路径必填，须在 `test-project/` 下） |

### Git 分支管理
| 命令 | 说明 |
|------|------|
| `/git-list-branch-all` | 列出所有未 accept 的文档管理分支 |
| `/git-list-branch-bug` | 列出 Bug 类型未 accept 分支 |
| `/git-list-branch-rp` | 列出需求类型未 accept 分支 |
| `/git-current-branch-info` | 打印当前分支完整信息（含 worktree 绑定） |
| `/git-current-branch-commit` | 提交当前分支变更 |
| `/git-current-branch-merge` | 合并当前分支到主分支 |
| `/git-current-branch-push` | 推送当前分支到远端 |
| `/git-init` | 初始化 Git 仓库 |
| `/git-load [远端]` | 从远端同步项目（本地为空时） |
| `/git-current-branch-delete` | 删除分支并移除绑定 worktree（main 禁止删除，需二次确认） |

### 体系帮助
| 命令 | 说明 |
|------|------|
| `/help` | Harness 体系介绍和使用指南 |
| `/gogogo` | Agent 异常中断后接续执行 |

---

## 推荐工作流

```
/design → /plan → /implement → /review → /accept
  │                    │            │
  └─ 可视化确认                      └─ /debug（如有问题）
```

---

## 约束条件

1. 绝不修改 `C:\Users\orion\.cursor\plugins\local\superpowers\` 下的任何文件
2. 仅操作 `d:\code\cursor\my_harness\` 范围内的内容
3. 所有 Rules 使用 `.mdc` 格式，存放于 `.cursor/rules/`
4. 所有 Skills 使用 `SKILL.md` 格式，存放于 `.cursor/skills/`
5. 所有 Commands 使用 `.md` 格式，存放于 `.cursor/commands/`
6. 所有 Subagents 使用 `.md` 格式，存放于 `.cursor/subagents/`
7. 所有 Hooks 使用 `hooks.json` + 脚本，存放于 `.cursor/hooks/`
