# /git-delete

删除当前分支及其绑定的 worktree。

## Agent 行为

1. 获取当前分支名
2. **禁止删除主分支**：若当前分支为 `main` 或 `master`，拒绝执行并提示「主分支不可删除」
3. 用 `git worktree list` 解析当前分支绑定的 worktree 路径（可能即当前目录）
4. 打印当前分支完整信息（参考 `/git-info` 的输出格式，含 worktree 绑定状态）
5. 使用 **AskQuestion** 工具要求用户二次确认删除，提示将删除分支、移除 worktree 目录及未推送提交风险
6. 用户确认后，按顺序执行（须在**主检出**所在目录操作，或先 `cd` 到主检出）：
   1. 若当前 shell 位于待删 worktree 内，先切换到主检出目录（`main` / `master` 所在路径）
   2. `git worktree remove <worktree-path>`（工作区有未提交变更时须先处理或经用户确认后加 `--force`）
   3. `git branch -d <branch>`（安全删除，仅删除已合并分支）
7. 如步骤 3 失败（分支未合并），提示用户可用 `git branch -D <branch>` 强制删除，并二次确认风险；worktree 已移除则仅补做分支删除
8. 确认删除成功并报告结果（分支名、已移除的 worktree 路径）

**说明**：删除顺序为先移除 worktree、再删分支，与 `doc-project-structure`「分支与 worktree 开发」一致。合入后清理亦应移除对应 worktree。
