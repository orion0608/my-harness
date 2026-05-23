# /accept

验收当前分支工作，准备合入。

## Agent 行为

1. 执行 Superpowers verification-before-completion 流程
2. 执行 `/docs-project-check`，确保分支文档完整
3. 读取 `doc-project-structure` 规则，执行分支合入合成流程
4. 执行 Superpowers finishing-a-development-branch 流程
