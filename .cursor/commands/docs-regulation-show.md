# /docs-regulation-show

查看项目文档格式规范。

## Agent 行为

1. 从已加载的 `doc-project-structure.mdc` 规则中提取文档清单、触发时机（见 Part C "分支过程文档" 和 Part D "全局跟踪文档"）和决策矩阵（见 Part E）
2. 若上下文中未加载该规则，则读取 `.cursor/rules/doc-project-structure.mdc`
3. 汇总输出：文档清单（必写/选写）、创建/触发时机；正文结构与写作要求不在规则中列举，统一指向文首映射的 `doc-file-definition-*` 技能
4. 若用户指定某个文档（或未指定但需查看格式），读取对应 `doc-file-definition-<文档名>` 技能，输出完整模板与编写规范
