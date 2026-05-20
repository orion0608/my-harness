# /git-current-branch-info

打印当前分支的完整信息。

## Agent 行为

1. 获取当前分支名
2. 输出以下信息（能获取到的全部输出）：

```
分支名称：<branch-name>
分支类型：<根据命名前缀推断：Bug / 需求 / 其他>
关联 worktree：<路径（如有）>
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

3. 分支类型推断规则：
   - `bug/`、`bugfix/`、`fix/` → Bug
   - `rp/`、`feat/`、`feature/` → 需求
   - 其他 → 未分类
