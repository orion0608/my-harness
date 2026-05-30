---
name: cursor-ide-browser-usage
description: >
  cursor-ide-browser MCP 通用工具序：会话序、lock、snapshot、交互、截图与证据要求。
  凡涉及网页查看、操作、验证且已获得 http(s):// URL 时 Read 本技能；场景专用扩展（CDP、P/I、报告）
  由引用方规则补充，本技能不引用其他技能。
---

# cursor-ide-browser-usage

Harness 内 **`cursor-ide-browser`** 的**通用 MCP 操作信源**。何时必须用浏览器、第三级门禁见 **`priority-resolution`** 规则；本技能只定义拿到 URL 之后如何调用 MCP。

## 调用前

1. 通过 **`CallMcpTool`** 调用 `cursor-ide-browser`（调用前查阅各工具 schema）
2. 目标 URL 须为 **`http(s)://`**；**禁止** `file://`
3. **禁止**用 curl / wget / 仅读 HTML 文件代替浏览器验证
4. 未 `browser_navigate` 且未 `browser_snapshot`（或截图）时，不得断言「页面正常」「测试通过」「按钮可用」

## 会话序（每回合一次）

| 顺序 | 动作 |
|------|------|
| 1 | `browser_tabs` — `action`: `list` |
| 2 | `browser_navigate` 到目标 URL |
| 3 | 路径逐步操作（见下节；可多条路径循环） |
| 4 | `browser_lock` — `action`: `unlock`（本回合曾 lock 时） |

## lock 顺序

- 先 `browser_navigate`，再 `browser_lock` — `action`: `lock`（**将交互**且未 lock 时）
- 纯浏览、无交互的单 URL 可跳过 lock

## 路径逐步操作

对每条用户路径循环；**禁止跳步**、**禁止无 snapshot 即 click**。

| 步序 | 工具 | 要求 |
|------|------|------|
| 1 | `browser_navigate` | 路径或新 URL 起点 |
| 2 | `browser_lock` | `action`: `lock`（将交互且未 lock 时；纯浏览可跳过） |
| 3 | `browser_snapshot` | 取 `ref`；`browser_click` / `browser_fill` / `browser_select_option` / `browser_type` **必须**用本步 `ref` |
| 4 | 交互（按需） | Tab / 菜单 / 填表 / 提交；**新 URL** → 回到步序 1 |
| 5 | `browser_snapshot` | 步序 4 之后或 DOM 可能变化时**必做** |
| 6 | `browser_take_screenshot` | 每 URL 至少一次；有交互则在里程碑再截 |
| 7 | 重复 3→6 | 直至该路径结束 |

**仅浏览、无交互的单 URL**：1 → 3 → 6（跳过 2、4、5）。

## 交互通则

- 任何可能改变页面结构的操作**前/后**：以**最新** `browser_snapshot` 为结构依据；交互后结构不明则再 snapshot
- 宣称完成或向用户汇报 Web 验证结果前：至少一次 snapshot 或 `browser_take_screenshot` 作为证据摘要
- 禁止在回复中长篇描述页面结构或交互，却不调用浏览器工具

## 适用场景（工具序相同）

| 场景 | 说明 |
|------|------|
| UI / HTML / Demo 预览与确认 | navigate + snapshot 确认布局与内容 |
| 前端功能验证、回归、手工测试 | 走关键用户路径，变更后重新 snapshot |
| Bug 复现与修复确认（Web） | 按复现步骤执行，修复后同路径再验证 |
| 设计 / 验收中的「能打开、能操作」 | 获得可访问 URL 后按本技能执行 |

**不适用**（无需强行开浏览器）：纯后端 / API / CLI / 单元测试且无 DOM 或浏览器行为；用户明确只要代码审查或静态分析。

## 场景扩展（本技能不定义）

引用方可在本工具序之上叠加专用步骤（如 Code Review 的 CDP 诊断、P/I 分级、报告模板），**不得**在引用方重复本节的逐步操作表。
