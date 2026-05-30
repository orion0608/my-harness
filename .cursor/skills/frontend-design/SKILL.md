---
name: frontend-design
description: >
  Single router for frontend UI: open aesthetics, ui-style DESIGN.md references,
  or ui-frame component libraries. Use for any UI layout, demo, or implementation.
---

# Frontend Design

前端 **style + frame + 动效 + 图表** 的唯一路由技能。`design.md`「前端设计」或用户明示为输入；何时 AskQuestion 由调用方（如 `brainstorming-visual-enhancement`）决定。

HTML 预览 → **`preview-html`**。

## 目录

| 路径 | 角色 |
|------|------|
| `free-style/SKILL.md` | 开放式美学子技能（无 `ui-frame` 库时加载） |
| `ui-style/<slug>/DESIGN.md` | 参考式数据（无独立子技能） |
| `ui-frame/<库>/SKILL.md` | 约束式子技能（如 `antd`） |
| `charts/echarts/SKILL.md` | 业务图表实现信源（ECharts 绑定与常见模式） |

## 模式与路由

| 美学 | 实现 | 行为 |
|------|------|------|
| 开放式 | 无库 | **Read** `free-style` 子技能 + 纯 HTML |
| 开放式 | antd 等 | 加载 `ui-frame` 子技能 |
| 参考式 | 任意 | **Read** `ui-style/<slug>/DESIGN.md`；有库则再加载对应 `ui-frame` 子技能 |
| （缺记录） | — | 不擅自默认；由调用方先确认 |

```
读 design.md「前端设计」或用户明示
→ 参考式 → 确认 slug（见下索引）→ Read ui-style/<slug>/DESIGN.md
→ 约束式 / 项目已有库 → 核验栈兼容（见下「ui-frame 与技术栈兼容性」）→ 加载 ui-frame 子技能（如 antd）
→ 开放式且无库 → **Read** free-style 子技能
→ 否则 → 纯 HTML/CSS/JS
→ 动效路由 → GSAP 插件技能（若命中，见下节）
→ 图表路由 → echarts 子技能（若命中，见下节）
```

**实现栈兜底**（无 design.md）：用户指定 > lockfile > 源码 import > 纯 HTML。

## 图表路由

业务数据可视化 **绑定与常见模式** 由 **`charts/echarts`** 子技能提供。本技能只定义 **何时 Read** 该子技能；禁止在 Harness 内重复 ECharts option 全文。

### 输入（按优先级）

1. `design.md`「前端设计」→ **图表库**
2. 用户明示（如「用 @ant-design/charts」）
3. 实现需求推断 + lockfile / 源码 import（见下「默认策略」）

### 判定

| 条件 | 行为 |
|------|------|
| 图表库 = `无` 或 `CSS only`，且无业务可视化需求 | 不上图表库；装饰性图形用 CSS/SVG |
| 图表库 = `ECharts`，或业务 dashboard / 趋势 / 占比 / 统计图 | **Read** `charts/echarts` 子技能 |
| 设计文档 `FLOW` / `ARCH` / `DATA` 结构图 | **Mermaid**（`brainstorming-visual-enhancement`）；不用 ECharts 替代 |
| lockfile 或源码已用 `@ant-design/charts` / G2 系 | 尊重既有栈，不强行换 ECharts |
| 用户指定 Chart.js、Recharts、D3 等 | 尊重选型，不强行 ECharts |

### 默认策略（design.md 未记录时）

- **含 dashboard、统计、趋势、占比、KPI 等业务图表** → **优先 ECharts**（非强制）；写代码前 **Read** `charts/echarts`
- **纯 CRUD、无图表** → 不上库
- **antd 后台且 lockfile 已有 @ant-design/charts** → 沿用 G2 系，不默认 ECharts
- **极简 sparkline / 纯装饰** → CSS/SVG 即可；不为微图形单独上 ECharts

### 次序

style/frame 路由完成 → 动效路由 → **图表路由**（写图表代码前）→ 按需 Read `charts/echarts` 子技能。

### 禁止

- 禁止未 Read `charts/echarts` 就写复杂 ECharts option（多 series、地图、联动等）
- 禁止在 Harness 技能内复制 ECharts API / option 全文
- 禁止用 Mermaid 或静态图替代产品 UI 中的交互式业务图表

## 动效路由

动效 **API 与最佳实践** 由 **GSAP 插件**（Cursor 外部插件）提供。本技能只定义 **何时 Read** 对应插件技能；禁止在 Harness 内重复 GSAP API 正文。

**前置**：须已安装 GSAP 插件（见 README「推荐外部插件」）。未安装则提示安装，不得在 Harness 内替代。

### 输入（按优先级）

1. `design.md`「前端设计」→ **动效强度** / **动画库**
2. 用户明示（如「用 GSAP ScrollTrigger」）
3. 实现需求推断（见下「默认策略」）

### 判定

| 条件 | 行为 |
|------|------|
| 动画库 = `无` 或 `CSS only`，或动效强度 = `静态` | 不读 GSAP 插件；CSS transition / `@keyframes` / `animation-timeline` |
| 动画库 = `GSAP`，或需 timeline 编排 | Read GSAP 插件（按场景选 core / timeline 等子技能） |
| 滚动联动 / pin / scrub / parallax | 同上，含 scrolltrigger；复杂场景含 performance |
| 需 Flip / Draggable 等插件能力 | 含 plugins 子技能 |
| React / Next 实现 | 含 react 子技能 |
| Vue / Svelte 等 | 含 frameworks 子技能 |

### 默认策略（design.md 未记录时）

- **ui-frame / antd 后台**：默认 `CSS only`；除非用户明确要求 GSAP。
- **纯 HTML Demo / 开放式 landing**：计划含 scroll 叙事、pin、stagger 入场 → 默认 **GSAP**。
- **简单 hover / opacity**：CSS 即可；不为微交互单独上 GSAP。
- 用户已选定其他动画库 → 尊重选型，不强行 GSAP。

### 次序

style/frame 路由完成 → **动效路由**（写动画代码前）→ 按需 Read GSAP 插件技能。

### 禁止

- 禁止 `window.addEventListener('scroll', …)` 做动画
- 禁止未 Read GSAP 插件技能就写 ScrollTrigger / timeline 代码
- 禁止在 Harness 技能内复制 GSAP API 正文

### 参考式读取

1. slug 须在下表索引中
2. 单任务通常只读 **1** 份 `DESIGN.md`；禁止扫读整个 `ui-style/`
3. 按色板、字体、组件态、Do/Don't 输出；与 antd 并用时做近似 token 映射
4. 可选 `preview.html` → 调用方执行 **`preview-html`**

### 开放式子技能

| 子技能 | 路径 |
|--------|------|
| `free-style` | `free-style/` |

美学为开放式且无 `ui-frame` 库时：**Read** `free-style` 子技能，再按项目栈实现（通常纯 HTML/CSS/JS）；Demo 须遵循 brainstorming 规范。

### 约束式子技能

| 子技能 | 路径 |
|--------|------|
| `antd` | `ui-frame/antd/` |

新增库：`ui-frame/<名>/SKILL.md` + 上表 + README。

#### ui-frame 与技术栈兼容性（必选检查）

加载任一 `ui-frame` 子技能前，**须确认该组件库支持当前项目的前端技术栈**（框架、运行时、构建目标）。不兼容时**不得**加载对应子技能；改选与栈匹配的库，或由调用方确认后回退 `free-style` / 纯 HTML。

依据（按优先级）：`design.md`「前端设计」已记录的栈 > 项目 lockfile / 配置文件 > 源码中的框架 import 与目录结构。

| 库 | 支持的技术栈 |
|----|-------------|
| `antd` | React + Web DOM；SSR/Electron；现代浏览器；TypeScript；npm/pnpm/yarn/bun；v6 需 React 18+ |

新增 `ui-frame` 子技能时，须在子技能 `SKILL.md` 或上表中写明**支持的技术栈**，避免 Agent 在错误栈上套用组件库。

## 参考索引（slug）

来源 awesome-design-md，扁平存放于 `ui-style/<slug>/`。

### AI & LLM

claude, cohere, elevenlabs, minimax, mistral.ai, ollama, opencode.ai, replicate, runwayml, together.ai, voltagent, x.ai

### 开发者工具

cursor, expo, lovable, raycast, superhuman, vercel, warp

### 后端 / DevOps

clickhouse, composio, hashicorp, mongodb, posthog, sanity, sentry, supabase

### 生产力 SaaS

cal, intercom, linear.app, mintlify, notion, resend, zapier

### 设计 / 创意

airtable, clay, figma, framer, miro, webflow

### 金融科技

binance, coinbase, kraken, mastercard, revolut, stripe, wise

### 电商 / 零售

airbnb, meta, nike, shopify, starbucks

### 媒体 / 消费科技

apple, ibm, nvidia, pinterest, playstation, spacex, spotify, theverge, uber, vodafone, wired

### 汽车

bmw, bmw-m, bugatti, ferrari, lamborghini, renault, tesla

### 其它

slack

### 图表子技能

| 子技能 | 路径 |
|--------|------|
| `echarts` | `charts/echarts/` |

新增图表库：`charts/<名>/SKILL.md` + 上表 + README + 图表路由判定表。

## 扩展

- 新参考：新增 `ui-style/<slug>/` + 更新本节索引
- 新库：`ui-frame/<名>/` + 约束式子技能表
- 新图表库：`charts/<名>/` + 图表子技能表 + 图表路由判定
- 开放式引导变更：更新 `free-style` 子技能（父技能表无需改，除非更名）
