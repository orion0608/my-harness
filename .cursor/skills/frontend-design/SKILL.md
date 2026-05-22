---
name: frontend-design
description: >
  Single router for frontend UI: open aesthetics, ui-style DESIGN.md references,
  or ui-frame component libraries. Use for any UI layout, demo, or implementation.
---

# Frontend Design

前端 **style + frame** 的唯一路由技能。`design.md`「前端设计」或用户明示为输入；何时 AskQuestion 由调用方（如 `brainstorming-visual-enhancement`）决定。

HTML 预览 → **`preview-html`**。

## 目录

| 路径 | 角色 |
|------|------|
| `free-style/SKILL.md` | 开放式美学子技能（无 `ui-frame` 库时加载） |
| `ui-style/<slug>/DESIGN.md` | 参考式数据（无独立子技能） |
| `ui-frame/<库>/SKILL.md` | 约束式子技能（如 `antd`） |

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
```

**实现栈兜底**（无 design.md）：用户指定 > lockfile > 源码 import > 纯 HTML。

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

## 扩展

- 新参考：新增 `ui-style/<slug>/` + 更新本节索引
- 新库：`ui-frame/<名>/` + 约束式子技能表
- 开放式引导变更：更新 `free-style` 子技能（父技能表无需改，除非更名）
