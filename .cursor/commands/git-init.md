# /git-init

将当前项目初始化为 Git 仓库。

## Agent 行为

1. 检查当前目录是否已是 Git 仓库
2. 如已是，报告当前状态并结束
3. 如不是，执行 `git init`
4. 按 doc-project-structure 规则创建 `.gitignore`（添加 `docs/*/`）
5. 执行初始提交（如项目已有文件）
6. 报告初始化结果
