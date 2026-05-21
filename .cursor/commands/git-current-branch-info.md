# /git-current-branch-info

打印当前分支的完整信息。

## Agent 行为

1. 获取当前分支名
2. 检测 worktree 绑定状态（见下节「worktree 检测」）
3. 输出以下信息（能获取到的全部输出）：

```
分支名称：<branch-name>
分支类型：<根据命名前缀推断：Bug / 需求 / 其他>
关联 worktree：<路径（如有）>
worktree 绑定：✅ 已绑定 / ⚠️ 主检出（main 正常）/ ❌ 未绑定（功能分支须在独立 worktree 中开发）
最后提交：<hash> - <message>（<时间>）
未推送提交数：<count>
落后/领先远端：<behind> / <ahead>
文档目录：✅ 存在 / ❌ 不存在

文档状态（如有 docs/<branch>/）：
  requirement.md   ✅ / ❌
  design.md        ✅ / ❌
  plan.md          ✅ / ❌
  devlog.md        ✅ / ❌
  accept-log.md    ✅ / ❌
  阶段推断：<design / plan / implement / review / accept>
```

4. 分支类型推断规则：
   - `bug/`、`bugfix/`、`fix/` → Bug
   - `rp/`、`feat/`、`feature/` → 需求
   - 其他 → 未分类

5. 若 `worktree 绑定` 为 ❌，追加一行提示：按 `doc-project-structure` 为当前分支创建 worktree（如 `git worktree add .worktrees/<branch-slug> -b <branch>` 或 `git worktree add .worktrees/<branch-slug> <branch>`），勿在主检出内仅 `checkout` 后继续开发。

### worktree 检测

在仓库根或当前目录执行：

```bash
git rev-parse --show-superproject-working-tree 2>/dev/null   # 非空 → 子模块，不按 worktree 判定
git rev-parse --git-dir
git rev-parse --git-common-dir
git worktree list --porcelain
```

判定（与 `doc-project-structure` 一致）：

| 条件 | `worktree 绑定` |
|------|-----------------|
| 当前分支为 `main` 或 `master`，且 `git-dir` 与 `git-common-dir` **相同** | ⚠️ 主检出（main 正常） |
| `git-dir` 与 `git-common-dir` **不同**（且非子模块） | ✅ 已绑定 |
| 功能分支，但 `git-dir` 与 `git-common-dir` **相同** | ❌ 未绑定 |

`关联 worktree`：从 `git worktree list` 中取当前分支对应条目的路径；若无则填「无」。
