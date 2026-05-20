---
name: doc-file-definition-devlog
description: 定义 devlog.md（开发日志）的编写格式和内容要求。当 doc-project-structure 规则触发 devlog.md 追加时使用。
disable-model-invocation: true
---

# devlog.md 编写规范

## 创建时机

每次 `git commit` 后追加一条记录。文件在第一次 commit 时创建。

## 文件位置

`docs/<branch-id>/devlog.md`

## 模板

```markdown
# 开发日志

> 分支：`<branch-name>`

---

## YYYY-MM-DD HH:MM

**Commit：** `<commit-hash>`

**摘要：** <!-- 一句话描述做了什么 -->

**详情：**
- <!-- 具体变更点 1 -->
- <!-- 具体变更点 2 -->

**影响范围：** <!-- 受影响的功能/模块 -->

---
```

## 内容要求

### Commit

- 使用完整的 commit hash（前 7 位即可）

### 摘要

- 一句话，让人秒懂这次提交做了什么
- 格式：`[类型] 描述`（如 `[feat] 添加用户登录`、`[fix] 修复日期格式错误`）

### 详情

- 列出具体的变更点，每个一行
- 如果是实现某个 Task，标注 Task 编号
- 如果修复了 issue，标注 issue 编号

### 影响范围

- 说明哪些模块或功能受到了影响
- 如果有破坏性变更，必须加粗标注

## 追加规则

- 每次 commit 后，在文件末尾追加一条新记录
- 记录之间用 `---` 分隔
- 时间使用提交时间（非写入时间）
