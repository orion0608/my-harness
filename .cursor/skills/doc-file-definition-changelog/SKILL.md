---
name: doc-file-definition-changelog
description: 定义 CHANGELOG.md（更新日志）的编写格式和内容要求。当 doc-project-structure 规则触发 CHANGELOG.md 创建或更新时使用。
disable-model-invocation: true
---

# CHANGELOG.md 编写规范

## 创建时机

项目初始化时创建。每次分支合并/发布时追加新版本的条目。

## 文件位置

`CHANGELOG.md`（项目根目录）

## 格式遵循

遵循 [Keep a Changelog](https://keepachangelog.com/) 规范，使用语义化版本。

## 模板

```markdown
# Changelog

本项目所有重要变更均记录于此。

格式遵循 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，
版本号遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [Unreleased]

<!-- 尚未发布的变更，按类别列出 -->

### Added
- 新增 xxx 功能

### Changed
- 将 xxx 从 A 迁移到 B

### Deprecated
- xxx 即将在下一大版本中移除

### Removed
- 移除了 xxx

### Fixed
- 修复了 xxx 在 yyy 场景下的 bug

### Security
- 修复了 xxx 安全漏洞（CVE-xxxx-xxxx）

## [1.0.0] - 2026-01-15

<!-- 首个正式发布版本 -->

### Added
- 项目的首个稳定版本
- 支持 xxx 核心功能

## [0.1.0] - 2025-12-01

<!-- 首个可用版本 -->

### Added
- 项目骨架搭建
- xxx 基础模块

---

<!-- 版本比较链接（可选）
[Unreleased]: https://github.com/user/repo/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/user/repo/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/user/repo/releases/tag/v0.1.0
-->
```

## 变更分类

| 类别 | 适用场景 |
|------|---------|
| `Added` | 新增功能 |
| `Changed` | 现有功能的变更 |
| `Deprecated` | 即将移除的功能 |
| `Removed` | 已移除的功能 |
| `Fixed` | Bug 修复 |
| `Security` | 安全修复 |

## 条目写法

- 每条以 `- ` 开头，用一句话描述变更
- 面向用户描述：用户能看到/感受到什么变化
  - ✅ "新增深色模式，可在设置中切换"
  - ❌ "添加了 ThemeContext 和 useTheme hook"
- 如涉及破坏性变更，用 **粗体** 标注

## 版本号规则

- 主版本号（MAJOR）：不兼容的 API 修改
- 次版本号（MINOR）：向下兼容的功能新增
- 修订号（PATCH）：向下兼容的问题修正
- `[Unreleased]` 段落始终在文件顶部，发布时改为具体版本号并添加日期

## 追加流程

分支合入时，从 `docs/<branch>/release-notes.md` 提取摘要，追加入 `[Unreleased]` 的对应分类下。发布时，将 `[Unreleased]` 下内容移入新版本段落。
