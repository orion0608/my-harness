# /review

对当前分支改动进行代码审查。

## Agent 行为

1. 确定审查范围（当前分支相对于主分支的 diff）
2. 执行 Superpowers requesting-code-review 流程
3. 审查完成后，按 doc-project-structure 决策矩阵记录审查结果
4. 有问题 → 回到 `/implement`；通过 → 提示使用 `/accept 命令`

