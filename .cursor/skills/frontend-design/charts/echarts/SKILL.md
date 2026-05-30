---
name: echarts-charting
description: >
  Business data charts with Apache ECharts (preferred default, not mandatory).
  Use for dashboards, trends, comparisons, composition, KPIs; framework bindings,
  resize/dispose, and common patterns. Read when frontend-design chart routing hits ECharts.
---

# ECharts Charting

> 位于 `charts/echarts/`；由 `frontend-design`「图表路由」在业务可视化且选型为 ECharts（或默认优先 ECharts）时 **Read** 加载。

**ECharts 为 Harness 业务图表的优先默认，非绝对强制**——用户、`design.md`「图表库」或 lockfile 已锁定其他库时，尊重既有选型。

## 适用 / 不适用

| 适用 | 不适用 |
|------|--------|
| line / bar / pie / scatter 及 dashboard 组合 | 设计阶段 FLOW/ARCH/DATA 结构图 → Mermaid |
| 趋势、对比、占比、KPI、多 series 业务图 | 纯 CSS/SVG 装饰、极简 sparkline |
| 纯 HTML Demo、React、Vue 产品页 | lockfile 已用 @ant-design/charts 且未要求迁移 |
| | 用户明确指定 Chart.js / Recharts / D3 等 |

## 框架绑定

写代码前确认项目栈，选用对应绑定：

| 栈 | 绑定 | 要点 |
|----|------|------|
| React 18+ | `echarts-for-react` 或 `useRef` + `echarts.init` | `useEffect` 内 init；cleanup 中 `chart.dispose()`；`ResizeObserver` 或 `window.resize` 调 `chart.resize()` |
| Vue 3 | `vue-echarts` | 组件卸载时 dispose；容器尺寸变化时 resize |
| 纯 HTML / Demo | CDN `echarts.min.js` + `echarts.init(dom)` | 容器须有明确 width/height；预览走 **`preview-html`** |

**版本**：优先与 lockfile 一致；无 lockfile 时用 ECharts 5.x 稳定 API。禁止混用已废弃的 v4 写法。

## 写代码前（必填上下文）

向实现中注入或自行确认：

1. **图表类型**（line / bar / pie / scatter / …）
2. **数据字段** + **一行示例**（如 `{ date: "2026-01-01", sales: 120 }`）
3. **框架**（React / Vue / 纯 HTML）
4. **空数据 / loading** 态如何处理

`series.data` 格式须与 `xAxis.type`（category / value / time）一致；多 series 时 legend 与 series.name 对应。

## 主题与色系

- 对齐 `design.md`「前端设计」色系或 antd Design Token
- 禁止硬编码与项目冲突的默认色（如无关的紫色渐变）
- 暗色模式项目：background、axisLabel、splitLine 须可读

## 常见模式（骨架，非 API 大全）

### 单折线（类目轴 + 数值）

- `xAxis: { type: 'category', data: [...] }`
- `yAxis: { type: 'value' }`
- `series: [{ type: 'line', data: [...] }]`
- `tooltip: { trigger: 'axis' }`

### 柱状对比

- 同上；`series.type: 'bar'`；多组对比用多 series 或 `dataset` + `encode`

### 饼图占比

- `series: [{ type: 'pie', radius: '60%', data: [{ name, value }, ...] }]`
- `tooltip: { trigger: 'item' }`

### 空数据 / loading

- 无数据：显示 Empty 占位或「暂无数据」文案，不要留空白 canvas
- 异步：先 loading 态，数据到达后 `setOption` 或更新 props

## 官方参考（按需查阅，不在 Harness 复制）

- [Handbook 入门](https://echarts.apache.org/handbook/zh/get-started.html)
- [Option 配置项](https://echarts.apache.org/zh/option.html)

复杂图（map、sankey、candlestick、3D）先查 Option 文档对应 series 类型，再写 option。

## 禁止

- 禁止编造已废弃或错误的 option 字段
- 禁止 React/Vue 中忽略 `dispose` / `resize` 导致泄漏或布局错乱
- 禁止在已有 @ant-design/charts 栈中擅自替换为 ECharts（除非用户或 design 已变更图表库）
- 禁止在本技能内展开完整 option 文档——细节以官方 Option 文档为准
