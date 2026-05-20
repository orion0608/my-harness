# /git-load [远端地址]

将远端项目同步到本地（仅当本地项目为空时可用）。

## Agent 行为

1. 检查本地项目根目录是否为空（仅 .cursor 等配置目录除外）
2. 如不为空，拒绝执行并提示原因
3. 确认远端地址有效
4. 执行 `git clone` 或 `git remote add origin` + `git pull`
5. 报告同步结果和项目结构
