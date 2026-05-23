# /proc-shutdown

关闭当前逻辑工作区所在项目、当前 Git 分支下的运行中服务实例。

## 前置

1. 确认当前逻辑工作区（`test-project/<name>/` 子项目或用户指定目标；Shell 操作以该路径为 cwd）
2. Read `tech-stack-guidance`；若项目为**工具级**且后端已集成 `go-lifecycle`，走注册表链路；否则按项目 README / 启动方式识别进程并关闭

脚本路径均相对于 **my_harness 项目根**：`.cursor/skills/tech-stack-guidance/go-lifecycle/scripts/`。

## Agent 行为（工具级 · go-lifecycle）

1. 解析 `AppName`（与 Go 库一致：`git remote get-url origin` 仓库名 → 失败则 `git rev-parse --show-toplevel` 目录名 → 再失败则 cwd 目录名）
2. 解析 `BranchName`（`git rev-parse --abbrev-ref HEAD`；detached HEAD 为 `detached-<short>`）
3. 查询注册表：

```bash
node .cursor/skills/tech-stack-guidance/go-lifecycle/scripts/svc-registry-branch.cjs --app <AppName> --branch <BranchName>
```

（在 my_harness 根目录执行，或通过绝对路径调用上述脚本）

4. 对 `files[].registry.instances[]` 中每条实例：
   - `svc-info.cjs --port <port>` 核实存活
   - 标注 `keepaliveStale: true` 的实例说明可能已异常退出，仍列入清单
5. **列出清单**（port、pid、workingDirectory、version、keepaliveStale），**AskQuestion** 确认要关闭的实例
6. 对确认项依次：`svc-shutdown.cjs --port <port>`
7. 再次执行步骤 3 或 `svc-info` 验证已关闭

## Agent 行为（非工具级或未集成 lifecycle）

1. 按项目文档 / 技术栈识别相关进程（端口、进程名、容器等）
2. 列出清单，**AskQuestion** 确认
3. 确认后关闭并验证

## 工具与规则

- 关闭前须有清单与用户确认（`priority-resolution` · AskQuestion）
- 宣称「已关闭」前须 `verification-before-completion` 证据（注册表无对应实例 / info 不可达 / 进程不存在）
