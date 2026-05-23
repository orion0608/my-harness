# /proc-list

列出当前逻辑工作区所在项目下**所有分支**正在运行的服务实例信息。

## 前置

1. 确认当前逻辑工作区（Shell 操作以该路径为 cwd）
2. Read `tech-stack-guidance`；**工具级**且已集成 `go-lifecycle` 时走注册表链路；否则按项目文档 / 进程工具列出同项目全部相关进程

脚本路径均相对于 **my_harness 项目根**：`.cursor/skills/tech-stack-guidance/go-lifecycle/scripts/`。

## Agent 行为（工具级 · go-lifecycle）

1. 解析 `AppName`（规则同 `proc-shutdown`）
2. 查询注册表（全分支）：

```bash
node .cursor/skills/tech-stack-guidance/go-lifecycle/scripts/svc-registry-project.cjs --app <AppName>
```

3. 对每条 `files[].registry.instances[]` 调用 `svc-info.cjs --port <port>`
4. 按 `branchName` 分组或表格汇报：`branchName`、`port`、`pid`、`workingDirectory`、`version`、`keepaliveStale`
5. 无实例时明确报告「项目无运行中服务」

## Agent 行为（非工具级或未集成 lifecycle）

1. 按项目约定识别同项目全部相关进程（含不同 worktree / 分支若可区分）
2. 列出可用信息；无法按分支区分时如实说明
