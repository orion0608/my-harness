# /implement

按实现计划开始编码。

## Agent 行为

1. 读取当前分支的计划文档（路径见 doc-project-structure）
2. 选择 Superpowers executing-plans 和 subagent-driven-development 执行，开展编程
3. 每次 commit 后、发生计划变更时，按 doc-project-structure 决策矩阵更新对应文档
4. 完成所有 Task 后提示下一步：`/review`

