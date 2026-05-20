# /debug [问题描述]

对指定问题启动系统化调试。

## Agent 行为

1. 记录问题（按 doc-project-structure 规则中 issues 的写入时机）
2. 执行 Superpowers systematic-debugging 流程
3. 修复完成后：
   - 按 doc-project-structure 决策矩阵更新相关文档
   - 对最后一次 commit diff 执行快速自审
