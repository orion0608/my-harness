# /myharness-test-project

声明在 `test-project/` 下开始针对某一子项目的 **Harness 验证**，并激活该子项目的**逻辑工作区**。

本命令**仅**负责：路径校验、（必要时）创建空目录、在会话中绑定 `<logical-root>`。开发流程、文档维护、Git 分支及后续步骤**一律**由 Harness 规则与其它命令在用户发起时接管——本命令不得探测子项目进度或编排后续流程。

物理 Cursor 工作区根目录必须保持为 `my_harness`（见 `my-harness-build` 规则「test-project」与「工作区约束」）。

## 用法

```
/myharness-test-project <project-name>
```


| 参数               | 必填  | 说明       |
| ---------------- | --- | -------- |
| `<project-name>` | ✅   | 子项目根目录名称 |


**合法路径示例**：

- `demo-app`（强制认为其物理路径为 ~/`test-project/demo-app` ，~代表my_harness项目的物理路径）

**非法示例**：

- `test-project/demo-app` — 禁止使用路径形式
- `../other` — 禁止路径穿越
- `D:\code\other-project` — 必须在 `test-project/` 下
- 空参数 — 须先指定路径

## Agent 行为

### 1. 解析并校验路径

1. 若未提供 `<project-name>`，列出 `test-project/` 下已有子目录（如有），并用 **AskQuestion** 或明确提示要求用户指定路径；**不得**自行猜测或默认选中。
2. 规范化路径：
  - 去除首尾空白与多余斜杠
  - 强制认为其物理路径为 ~/`test-project/<project-name>` ，~代表my_harness项目的物理路径
3. 解析后的**逻辑开发根目录**记为 `<logical-root>`（相对于 `my_harness` 项目根，如 `test-project/demo-app`）。

### 2. 检查子项目目录

1. 确认 `<logical-root>` 在磁盘上存在且为目录。
2. 若不存在：询问用户是否在 `<logical-root>` 创建空目录；仅在用户明确同意后才创建（**不**初始化 Git、不写文档、不搭骨架，除非用户另下指令）。
3. 若存在：继续下一步。

### 3. 激活逻辑工作区

在对话中**显式声明**并在此后全程遵守，直至用户切换子项目或明确结束验证：

- **Harness 根目录**：`my_harness` 项目根（Cursor 工作区根，不变）
- **逻辑开发根目录**：`<logical-root>`
- **开发目的**：验证 Harness
- **路径约定**：子项目的读写、Shell、Git 以 `<logical-root>` 为项目根；Harness 元数据（`.cursor/` 等）仍以 `my_harness` 根为基准

工作区约束与跨项目隔离：**仅引用并遵守** `my-harness-build` 规则（禁止物理切换工作区、禁止读取其它测试子项目等）。**不得**在本命令正文中复述该规则的完整表格。

### 4. 输出确认摘要

向用户输出固定格式的确认块，例如：

```
✅ 逻辑工作区已激活

Harness 根目录（物理工作区）：<my_harness 绝对路径>
逻辑开发根目录：           <logical-root 绝对路径>
开发目的：                 Harness 验证

约束：物理工作区不切换；跨项目隔离见 my-harness-build 规则。
后续：由用户自行发起 Harness 体系内的命令；本命令不预判流程阶段。
```

### 5. 禁止在本命令内执行

- ❌ 探测或汇报子项目 Git、分支、技术栈、分支过程文档目录等状态
- ❌ 根据状态建议或自动引导任何后续 Harness 命令或 Git 操作
- ❌ 编写或更新子项目的过程文档（文档维护见 `doc-project-structure` 规则，在用户触发对应业务时机时执行）

---

## 与其它组件的关系

- 激活后，用户在 Harness 体系内发起的、面向「当前项目」的操作，默认以 `<logical-root>` 为逻辑项目根（具体边界见 `my-harness-build` 规则）。
- Harness 元数据（`.cursor/` 等）的修改仍仅限 `my_harness` 可操作范围，见 `my-harness-build` 规则。

