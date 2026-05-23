# /git-list-branch-all

列出所有基于文档管理体系、尚未 accept 的分支。

## Agent 行为

1. 获取所有本地分支
2. 排除已合入主分支的分支（已 accept）
3. 对每个分支：若存在绑定 worktree，检查 `{worktree}/docs/<branch-name>/`；否则标记「无 worktree / 无过程文档目录」
4. 输出分支名、最后提交时间、文档目录状态
