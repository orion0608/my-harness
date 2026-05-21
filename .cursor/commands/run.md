# /run

启动项目开发服务。

## Agent 行为

1. 检测项目技术栈和依赖
2. 确认依赖已安装（未安装则提示）
3. 按依赖顺序启动所有服务
4. 验证启动成功并报告访问方式（须有证据，见下节）

## 工具与规则（启动与展示）

涉及前后端服务启动后的**访问展示**，或宣称「服务已就绪 / 页面可访问」时，须遵循下表；**不得**用 `curl`/`wget`、仅读 HTML 文件或 `file://` 代替浏览器验证。

| 场景 | 须遵循 | Agent 做法 |
|------|--------|------------|
| 有 http(s) 可访问的前端或 Web 页面 | `priority-resolution` | **cursor-ide-browser**：`browser_navigate` → 必要时 `browser_lock` → `browser_snapshot` 或 `browser_take_screenshot` 作为就绪证据；汇报时给出完整 `http(s)://` URL |
| 仅本地静态 HTML 目录、无 dev server | `preview-html` | 读取并执行该技能启动 HTTP 服务，用返回的 `url` 再 `browser_navigate` |
| 项目使用 Visual Companion 预览 | `brainstorming-visual-enhancement` | Windows 按该规则 Visual Companion 启动备忘；浏览器打开 Companion 的 `url`（禁止 `file://`） |
| 宣称启动或展示成功 | `verification-before-completion` | 终端日志或浏览器 snapshot/截图至少一项，再向用户汇报 |

**禁止**：未 `browser_navigate`（或等价打开）且未 snapshot/截图即断言「页面正常」「服务可用」。
