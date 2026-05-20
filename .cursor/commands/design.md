# /design [主题]

启动完整设计流程。

## Agent 行为

1. 确认当前分支（未创建则提醒）
2. 按 doc-project-structure 规则创建分支文档目录和初始文档
3. 执行 Superpowers brainstorming 流程
4. 涉及可视化触发条件时，遵循 brainstorming-visual-enhancement 规则；预览 `docs/<branch>/demo/` 下 HTML 时读取 `preview-html` 技能
5. 涉及架构决策触发条件时，遵循 brainstorming-architecture-thinking-enhancement 规则（在 "Present design" 前接入 architecture-design 技能）
6. 设计定稿后，按 doc-project-structure 决策矩阵写入设计文档
7. 提示下一步：`/plan`
