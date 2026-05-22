# /implement

按实现计划开始编码。

## Agent 行为

1. 读取当前分支的计划文档（路径见 doc-project-structure）
2. 选择 Superpowers executing-plans 和 subagent-driven-development 执行，开展编程
3. 完成所有 Task 后，读取并执行 `/review` 命令（视同用户已发出该指令；不仅提示下一步）

