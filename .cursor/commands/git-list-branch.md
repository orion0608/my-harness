# /git-list-branch [类型]

列出基于文档管理体系、尚未 accept 的分支。

**类型**（可选）：用户未指明时列出全部；用户声明 `bug` / `Bug` 或 `rp` / `需求` 等时按对应类型筛选。

| 用户意图 | 筛选规则 | 额外标注 |
|---------|---------|---------|
| 未指明 / 全部 | 所有未 accept 分支 | — |
| Bug | 分支名以 `bug/`、`bugfix/`、`fix/` 开头 | 未解决问题数（读 `{worktree}/docs/<branch-id>/issues.md`，无文件则填 0 或「—」） |
| 需求（rp） | 分支名以 `rp/`、`feat/`、`feature/` 开头 | 当前阶段（design / plan / implement / review / accept，推断规则同 `/git-info`） |

## Agent 行为

1. 解析用户是否声明类型：
   - 命令参数 `[类型]`，或用户在同一会话中明确说「Bug 分支」「需求分支」「rp」等
   - 无法判定时默认**全部**，不 AskQuestion
2. 获取所有本地分支，排除已合入主分支的分支（已 accept）
3. 按上表应用类型筛选（若已声明）
4. 对每个分支：
   - 用 `git worktree list` 解析绑定 worktree
   - 若存在 worktree，检查 `{worktree}/docs/<branch-id>/`（`<branch-id>` = 分支完整名）；否则标记「无 worktree / 无过程文档目录」
   - Bug 类型：统计 `issues.md` 中未关闭/未解决条目数（若无文件则 0）
   - 需求类型：根据过程文档存在性与内容推断阶段（与 `/git-info` 一致）
5. 输出：分支名、类型（Bug / 需求 / 其他）、最后提交时间、文档目录状态、以及类型相关的额外列（Bug：未解决问题数；需求：阶段）
