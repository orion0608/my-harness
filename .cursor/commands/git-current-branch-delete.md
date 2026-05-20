# /git-current-branch-delete

删除当前分支。

## Agent 行为

1. 获取当前分支名
2. **禁止删除主分支**：若当前分支为 `main` 或 `master`，拒绝执行并提示"主分支不可删除"
3. 打印当前分支完整信息（参考 `/git-current-branch-info` 的输出格式）：
   ```
   分支名称：<branch-name>
   分支类型：<根据命名前缀推断：Bug / 需求 / 其他>
   最后提交：<hash> - <message>（<时间>）
   未推送提交数：<count>
   落后/领先远端：<behind> / <ahead>
   ```
4. 使用 **AskQuestion** 工具要求用户二次确认删除，提示删除风险和影响
5. 用户确认后，先自动切换到 `main`（或 `master`）分支
6. 执行 `git branch -d <branch>`（安全删除，仅删除已合并分支）
7. 如安全删除失败（分支未合并），提示用户可用 `git branch -D <branch>` 强制删除，并二次确认风险
8. 确认删除成功并报告结果
