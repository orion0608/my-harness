---
name: tech-stack-guidance
description: >-
  组织固定技术栈与系统开发模式（工具级/平台级）的唯一约束信源。
  由 brainstorming-architecture-thinking-enhancement 在 0-1 架构接入时加载；
  用户选定级别，Agent 说明约束不替用户判定。
disable-model-invocation: true
---

# 技术栈开发指引

组织级固定技术栈与**工具级 / 平台级**开发约束的唯一信源。

---

## 何时加载

| 场景 | 行为 |
|------|------|
| **0-1**，且 `brainstorming-architecture-thinking-enhancement` 已触发架构接入 | **Read 本技能** → 说明两类定义与约束 → **AskQuestion** 用户选定 → 写 requirement §1 |
| **1-N** | **不** AskQuestion 重选；从项目 `README.md` / `docs/ARCHITECTURE.md` 读已锁定级别并校验实现 |
| 未触发架构接入 | 不加载选型流程 |

---

## 两类定义（给用户说明，Agent 不判定）

### 工具级（微型服务 / 工具级应用）

业务相对简单、以**本机运行**为主的软件或工具；典型为本地助手、桌面/移动工具、单机小业务。架构倾向单体（至多前后端两个程序），数据与配置在用户目录，**不使用**下文「云端组件」清单中的运行时依赖。

### 生产平台级

**其它**面向云端与生产环境的系统：微服务、容器部署、集中运维与可观测性。适用多用户业务平台、需 PostgreSQL/Redis/RabbitMQ 等云端组件的场景。

**选定由用户决定**；Agent 并列展示两类约束差异后 AskQuestion，**禁止**默认工具级或平台级。

---

## 全局技术栈（组织固定，不可替换）

| 层次 | 选型 |
|------|------|
| 后端 | Go（多平台） |
| Web / H5 | React + Vite |
| 桌面 | Windows / Mac：Capacitor（H5 混合） |
| 移动 | iOS / 鸿蒙 / 微信小程序：Capacitor（H5 混合）；安卓待项目补充 |
| 本地关系型库 | SQLite |
| 本地分析型库 | DuckDB（本地分析、CSV/Parquet） |
| 云端关系型库 | PostgreSQL |
| 云端 KV | Redis |
| 云端时序库 | InfluxDB |
| 云端消息队列 | RabbitMQ |
| 云端代理 | nginx |

**云端组件**指上表中 PostgreSQL、Redis、InfluxDB、RabbitMQ、nginx 及同类自建/托管中间件。**工具级**不得将其作为运行时依赖；调用公网 HTTP API（更新检查、第三方 REST）不视为引入云端组件。

---

## 工具级约束

| 项 | 要求 |
|----|------|
| 架构 | 单体；至多前端 + 后端两个程序 |
| 端口 | 前后端占用 **50000–60000** 间随机端口 |
| 标准接口 | 服务信息查询（启动时间、路径、版本、进程 ID 等）；服务关闭（收到通知后退出） |
| 排他启动 | 启动前检查旧实例；先调关闭接口，失败再用命令行结束进程 |
| 数据目录 | 用户路径下，如 Windows：`C:\Users\<用户名>\.<应用名>` |
| 配置 | JSON 本地文件 |
| 持久化 | 本地数据库（SQLite）；分析场景可用 DuckDB |
| 禁止 | 不得选用云端组件；逻辑尽量用语言自实现 |

---

## 平台级约束

| 项 | 要求 |
|----|------|
| 架构 | 微服务 + 容器部署；容器监管推荐 [Portainer](https://github.com/portainer/portainer) |
| 服务划分 | 按业务域包干，不宜过碎 |
| 接口 | 每服务三类：业务接口、对内服务接口、对外系统接口 |
| 鉴权 | 本地配置各服务访问秘钥；统一签名验签 |
| 数据 | 共用 DB 实例；各服务独立逻辑库；跨域数据须走对应服务接口 |
| 日志 | OpenTelemetry 生成 trace id；PLG：Promtail → Loki → Grafana |

---

## 用户选型流程

1. 并列简述两类定义（上一节）。
2. 对照**工具级约束**与**平台级约束**表，说明差异。
3. **AskQuestion**：`工具级` | `平台级`（禁止 Agent 代选）。
4. 将用户原话或选项写入 `docs/<branch>/requirement.md` **§1 用户原始需求**（不新增模板字段）。

---

## 级别锁定

- 项目**一旦选定级别即不可变更**。
- 禁止在同一项目内 in-place 改级（不得仅改文档字段继续开发）。
- 若需换级别：**新建项目**并重构；旧项目可归档。

**1-N**：不得 AskQuestion 重选；实现/审查只校验是否违反已锁定级别。

---

## 文档记录（不修改 requirement/design 模板结构）

| 文档 | 写入位置 | 内容 |
|------|----------|------|
| `requirement.md` | **§1 用户原始需求** | 用户原话或 AskQuestion 选项（含 工具级\|平台级） |
| `design.md` | **§5 关键决策** | 级别；与 requirement §1 一致 |
| `design.md` | **§6 风险与约束** | 级别锁定：本项目不可变更；改级须新建项目并重构 |
| `README.md` | **§技术栈** 表 | 「系统级别」行（格式见 `doc-file-definition-readme`） |
| `docs/ARCHITECTURE.md` | **§1 / §2 / §7** | 系统级别及锁定策略（格式见 `doc-file-definition-architecture`） |

---

## 自检清单

**工具级**

- [ ] 未引入 PostgreSQL / Redis / InfluxDB / RabbitMQ 等云端组件
- [ ] 单体或仅前后端；端口在 50000–60000
- [ ] 具备信息查询与关闭接口；排他启动
- [ ] 数据与 JSON 配置在用户目录；业务持久化用 SQLite

**平台级**

- [ ] 微服务边界与三类接口清晰
- [ ] 签名校验与逻辑库隔离
- [ ] OTel + PLG 日志方案已纳入设计/部署

**通用**

- [ ] 后端 Go；Web/H5 React + Vite；多端 Capacitor 按项目范围
- [ ] 文档链已按「文档记录」节落盘；README / ARCHITECTURE 在 0-1 定级后已归档
