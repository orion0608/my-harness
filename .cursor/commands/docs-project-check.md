# /docs-project-check

检查当前分支文档完整性。

## Agent 行为

1. 读取 doc-project-structure 规则，获取：
   - 分支过程文档清单（必写/选写）
   - 全局跟踪文档清单
2. 在**绑定 worktree** 内扫描 `docs/<branch-id>/`（路径解析见 `doc-project-structure`），逐项比对输出 ✅/❌
3. 输出缺失清单及建议补齐时机
