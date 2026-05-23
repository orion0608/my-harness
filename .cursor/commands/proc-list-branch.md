# /proc-list-branch

列出当前逻辑工作区所在项目、当前 Git 分支下正在运行的服务实例信息。

## 前置

1. 确认当前逻辑工作区（Shell 操作以该路径为 cwd）
2. Read `tech-stack-guidance`；**工具级**且已集成 `go-lifecycle` 时走注册表链路；否则按项目文档 / 进程工具列出同项目同分支进程

脚本路径均相对于 **my_harness 项目根**：`.cursor/skills/tech-stack-guidance/go-lifecycle/scripts/`。

## Agent 行为（工具级 · go-lifecycle）

1. 解析 `AppName`、`BranchName`（规则同 `proc-shutdown`）
2. 查询注册表：

```bash
node .cursor/skills/tech-stack-guidance/go-lifecycle/scripts/svc-registry-branch.cjs --app <AppName> --branch <BranchName>
```

3. 对 `files[].registry.instances[]` 中每条实例调用 `svc-info.cjs --port <port>` 获取实时 info
4. 以表格汇报：`port`、`pid`、`workingDirectory`、`version`、`startedAt`、`lastKeepalive`、`keepaliveStale`
5. 无实例时明确报告「当前分支无运行中服务」

## Agent 行为（非工具级或未集成 lifecycle）

1. 按项目约定识别同项目、同分支相关进程
2. 列出 pid、命令行、端口等可用信息；无法识别时说明依据与局限
