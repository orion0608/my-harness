# /git-current-branch-merge

将当前分支合并到主分支。

## Agent 行为

1. 检查当前分支是否已完成 accept 流程（文档完整性、测试通过）
2. 如未完成 accept，提示用户先执行 `/accept`
3. 如已完成，获取主分支名（main 或 master）
4. 切换到主分支并执行合并
5. 确认合并成功并报告结果
