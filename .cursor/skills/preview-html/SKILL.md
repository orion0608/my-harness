---
name: preview-html
description: >
  Serve local HTML folders over HTTP (zero deps) and open in Cursor's built-in
  browser. Supports multiple --root directories in parallel (one server each).
  Auto-reuse per directory.
---

# preview-html

把本地 HTML 目录用 HTTP 代理到 `http://localhost`，并在 **Cursor IDE 内置浏览器** 中打开。

**支持多目录并行**：一个后台 **daemon** 管理多个 `--root`（各用不同端口）。适合多个项目同时开发。

## 脚本

| 文件 | 路径（相对本技能目录） |
|------|------------------------|
| 静态服务器 | `scripts/demo-preview-server.cjs` |

## 快速使用

1. **定目录**：`--root` 指向包含 `.html` 的文件夹。
2. **确保服务运行**（每个目录首次需后台启动一次）：

```bash
node .cursor/skills/preview-html/scripts/demo-preview-server.cjs --root <目录A>
node .cursor/skills/preview-html/scripts/demo-preview-server.cjs --root <目录B>
```

3. **读 stdout**：首行 JSON。`reused: true` 表示该目录已有实例，无需新进程。
4. **打开预览**：**cursor-ide-browser** → `browser_navigate` 到 JSON 的 `url`。
5. **同目录再次预览**：再执行同 `--root` 即可复用；换页面只改 `browser_navigate`。
6. **结束**：`--stop --root <目录>` 只停该目录；`--stop-all` 停全部。

## 空闲与退出（daemon）

| 规则 | 时长 |
|------|------|
| 单个 `--root` 无 HTTP 请求 | **5 分钟**后自动关闭该 root |
| 该 root 收到任意请求 | 重新计时 5 分钟 |
| 全部 root 均已关闭 | 再等 **5 分钟**后 daemon 进程退出 |
| 手动 `--stop-all` | 立即关闭所有 root 并退出 daemon |

## 多目录并行

| 情况 | 行为 |
|------|------|
| 项目 A、B 各有一个 demo 目录 | 各 `--root` 一次，两个 `url`（不同端口），同一 daemon |
| 同一 `--root` 再次启动 | 复用（`reused: true`），并刷新该 root 的空闲计时 |
| 排他 | 同一目录只对应一个服务 |

状态：`%TEMP%/preview-html-daemon.json` + 每目录 `preview-html-<hash>.state.json`。

## 命令

| 命令 | 说明 |
|------|------|
| `--root <dir>` | 启动或复用该目录 |
| `--stop --root <dir>` | 仅停止该目录实例 |
| `--stop-all` | 停止所有实例 |
| `--status --root <dir>` | 查询该目录实例 |
| `--list` | 列出当前所有运行中实例（JSON） |
| `--port` / `--host` | 仅新启动时生效 |

**并行示例 stdout**：

```json
{"type":"preview-html-started","reused":false,"pid":1001,"port":52341,"url":"http://localhost:52341/","root":"D:/project-a/docs/demo/demo"}
{"type":"preview-html-started","reused":false,"pid":1002,"port":53412,"url":"http://localhost:53412/","root":"D:/project-b/docs/rp-x/demo"}
```

## 访问地址

| 目的 | URL |
|------|-----|
| 索引 | JSON 的 `url` |
| 某个页面 | `{url}<文件名>.html` |

## 在 Cursor 中打开

- ✅ **cursor-ide-browser** → `browser_navigate`
- ❌ `file://`

## 停止服务

- 单目录：`--stop --root <dir>`（立即关闭该 root）
- 全部：`--stop-all`（立即关闭 daemon）
- 或等待上述 5 分钟空闲规则自动回收
