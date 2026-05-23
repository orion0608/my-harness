---
name: visual-companion-windows
description: >
  Start Superpowers Visual Companion on Windows (Git Bash or WSL paths),
  write screen HTML, open in cursor-ide-browser. Use when brainstorming UI
  direction preview or Harness rules require Companion on win32.
---

# Visual Companion（Windows）

Superpowers `visual-companion.md` 以 macOS/Linux 为主。**凡使用 Visual Companion（含 `brainstorming-visual-enhancement`、`/run` 等）须先 Read 本技能**；Windows 启动与屏写以本技能为准。HTML 降级预览见 `preview-html`。

## 根因

| 现象 | 原因 |
|------|------|
| `bash ... start-server.sh` → 退出码 **127** | `C:\Windows\System32\bash.exe` 是 WSL 启动器，不是 Git Bash |
| `No such file or directory` | WSL 无法识别 `C:/Users/...` 路径 |
| 无 `.superpowers/brainstorm/<session>/` | 脚本未真正执行 |

**禁止**：未确认 shell 即用 `bash "C:/Users/.../start-server.sh"`。

## 启动（按优先级）

**方式 1 — Git Bash + `--background`（首选）**

```powershell
& "C:\Program Files\Git\bin\bash.exe" -lc `
  "'/c/Users/<user>/.cursor/plugins/local/superpowers/skills/brainstorming/scripts/start-server.sh' --project-dir '/d/code/cursor/my_harness/test-project/<name>' --background"
```

- `C:\` → `/c/`；`D:\` → `/d/`
- Shell **后台**执行；从 stdout 或 `<project>/.superpowers/brainstorm/<session>/state/server-info` 读取 `url`、`screen_dir`、`state_dir`

**方式 2 — WSL**

```powershell
wsl bash /mnt/c/Users/<user>/.cursor/plugins/local/superpowers/skills/brainstorming/scripts/start-server.sh `
  --project-dir /mnt/d/code/cursor/my_harness/test-project/<name> --background
```

- 必须用 `/mnt/c/...`，禁止 `C:/...`

## 预览页写法（经验）

Companion 内 **`<iframe>` 无法正常加载**（含 `src` 指向本地或其它页）。方向预览与屏上 mockup 须用 **单页内联**：布局、样式、脚本写在同一 HTML 片段内；多方案用 Tab/分段/并排块，**勿**用 iframe 嵌第二份 HTML 或外链页。

## 启动后

1. 确认 `state/server-info`；无或 `server-stopped` → 重试（空闲约 30 分钟自动退出）
2. **Write** 新 `*.html` 到 `screen_dir`（单页内联片段即可，见 Superpowers `visual-companion.md`；勿复用文件名）
3. **cursor-ide-browser** → `browser_navigate` 到 `url`（**禁止** `file://`）
4. 读 `state/events`（若有）合并用户反馈
5. 不需浏览器时推送 `waiting.html`；结束可跑 `stop-server.sh`（同样 Git Bash 或 WSL 路径）

## 失败降级

Git Bash 与 WSL 均失败：**告知原因** → `docs/<branch>/demo/direction-<简述>.html` → Read **`preview-html`** + 内置浏览器。**勿** AskQuestion 是否改 HTML；**勿**未尝试即降级或仍称已用 Companion。

## 本地 HTML（非 Companion）

`docs/<branch>/demo/` 下 HTML（阶段二 Demo 等）：**preview-html** + 内置浏览器（禁止 `file://`）。
