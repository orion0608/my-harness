# /git-current-branch-commit

提交当前分支的变更。

## Agent 行为

1. 显示当前未暂存和已暂存的变更摘要
2. 按 doc-project-structure 规则更新分支过程文档（devlog 追加）
3. 生成符合规范的 commit message（类型 + 简短描述）
4. 执行 `git add` 和 `git commit`
5. 确认提交成功并显示 commit hash
