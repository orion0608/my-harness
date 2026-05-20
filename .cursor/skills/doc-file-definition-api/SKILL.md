---
name: doc-file-definition-api
description: 定义 API.md（API 参考文档）的编写格式和内容要求。当 doc-project-structure 规则触发 API.md 创建或更新时使用。
disable-model-invocation: true
---

# API.md 编写规范

## 创建时机

首个 API 端点出现时创建。API 变更时更新。

## 文件位置

`docs/API.md`

## 模板

```markdown
# API 参考

> 最后更新：YYYY-MM-DD
> 基准 URL：`https://api.example.com/v1`

## 1. 概述

<!-- API 的整体说明：版本策略、内容类型、字符编码 -->

- 内容类型：`application/json; charset=utf-8`
- 认证方式：Bearer Token（详见 [认证](#认证)）
- 版本策略：URL 路径版本（`/v1/`、`/v2/`）

## 2. 认证

<!-- 认证方式说明 -->

所有需要认证的端点需在请求头携带：

```
Authorization: Bearer <token>
```

Token 通过 `POST /auth/login` 获取，有效期 24 小时。

## 3. 通用约定

### 错误模型

错误类型对齐 `architecture-design` 接口契约规定，分 5 类。每个端点的错误响应必须能映射到下表之一。

| 错误类型 | 含义 | 默认 HTTP 状态码 | 默认可重试 |
|---|---|---|:---:|
| `ValidationError` | 参数 / 格式 / 校验失败 | 400 / 422 | 否 |
| `AuthError` | 认证或授权失败 | 401 / 403 | 否 |
| `DomainError` | 业务规则违反（含资源冲突、状态不允许） | 409 / 422 | 否 |
| `InfrastructureError` | DB / 外部 SDK / 网络错误 | 502 / 503 / 504 | 是 |
| `SystemError` | 未分类错误 | 500 | 否 |

### 幂等性约定

| 方法 | 默认幂等性 | 说明 |
|---|:---:|---|
| GET / HEAD / OPTIONS | ✅ | 始终幂等 |
| PUT / DELETE | ✅ | 资源级幂等 |
| PATCH | ⚠️ | 需端点逐一标注 |
| POST | ❌ | 默认不幂等；需要幂等时使用 `Idempotency-Key` 请求头 |

每个 POST / PATCH 端点的「端点详情」必须显式标注幂等性。

### 版本与弃用

- 主版本通过 URL 路径切换：`/v1/`、`/v2/`
- 弃用端点必须在「端点详情」标注 `弃用` 状态、替代端点、计划下线时间
- 破坏性变更必须升主版本，不得在已发布版本内变更字段含义

### 分页

列表类接口统一使用以下查询参数：

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `page` | int | 1 | 页码（1-based） |
| `page_size` | int | 20 | 每页条数（最大 100） |

响应体包含分页元信息：

```json
{
  "data": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

### 错误响应

所有错误返回统一格式：

```json
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "用户可读的错误描述"
  }
}
```

### 时间格式

所有时间字段使用 ISO 8601 格式（UTC）：`2026-01-15T08:30:00Z`

## 4. 端点目录

<!-- 按资源分组列出所有端点 -->

### 认证

| 方法 | 路径 | 认证 | 说明 |
|------|------|:----:|------|
| POST | `/auth/login` | 否 | 用户登录 |
| POST | `/auth/refresh` | 否 | 刷新 Token |
| POST | `/auth/logout` | 是 | 登出 |

### 用户

| 方法 | 路径 | 认证 | 说明 |
|------|------|:----:|------|
| GET | `/users` | 是 | 获取用户列表 |
| GET | `/users/:id` | 是 | 获取用户详情 |
| POST | `/users` | 否 | 注册新用户 |
| PATCH | `/users/:id` | 是 | 更新用户信息 |
| DELETE | `/users/:id` | 是 | 删除用户 |

## 5. 端点详情

<!-- 每个端点展开描述 -->

### POST /auth/login

用户登录，返回 JWT Token。

**元信息：**

| 项 | 值 |
|---|---|
| 幂等性 | 否（同密码多次请求都会签发新 Token） |
| 幂等键 | 不适用 |
| 版本 | v1（当前） |
| 弃用计划 | 无 |

**请求体：**

```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:----:|------|
| `email` | string | 是 | 注册邮箱 |
| `password` | string | 是 | 密码（6-128 字符） |

**成功响应 `200`：**

```json
{
  "token": "eyJhbGciOi...",
  "expires_at": "2026-01-16T08:30:00Z",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "张三"
  }
}
```

**错误响应：**

| 状态码 | code | 错误类型 | 可重试 | 说明 |
|--------|------|---------|:----:|------|
| 400 | `VALIDATION_ERROR` | ValidationError | 否 | 请求参数不合法 |
| 401 | `INVALID_CREDENTIALS` | AuthError | 否 | 邮箱或密码错误 |
| 429 | `RATE_LIMITED` | InfrastructureError | 是 | 登录频率过高，建议指数退避重试 |

## 6. 错误码参考

| 状态码 | code | 错误类型 | 可重试 | 说明 |
|--------|------|---------|:----:|------|
| 400 | `VALIDATION_ERROR` | ValidationError | 否 | 请求参数校验失败 |
| 401 | `UNAUTHORIZED` | AuthError | 否 | 未认证或 Token 过期 |
| 403 | `FORBIDDEN` | AuthError | 否 | 无权限访问该资源 |
| 404 | `RESOURCE_NOT_FOUND` | DomainError | 否 | 请求的资源不存在 |
| 409 | `CONFLICT` | DomainError | 否 | 资源冲突（如重复创建） |
| 422 | `BUSINESS_ERROR` | DomainError | 否 | 业务规则不满足 |
| 429 | `RATE_LIMITED` | InfrastructureError | 是 | 请求频率超限 |
| 500 | `INTERNAL_ERROR` | SystemError | 否 | 服务器内部错误 |
| 502/503/504 | `UPSTREAM_ERROR` | InfrastructureError | 是 | 上游服务不可用 |
```

## 章节要求

### 概述

- 标注基准 URL、内容类型、认证方式
- 说明 API 版本策略（路径版本 / 请求头版本 / 查询参数版本）

### 认证

- 说清 Token 获取方式、携带方式、有效期
- 如有多套认证体系，分别说明

### 通用约定

- 必须涵盖：**错误模型**（对齐 `architecture-design` 5 类）、**幂等性约定**、**版本与弃用**、分页、错误响应格式、时间格式
- 错误模型 5 类（ValidationError / AuthError / DomainError / InfrastructureError / SystemError）必须显式声明
- 幂等性按 HTTP 方法默认值列出，POST / PATCH 由端点逐一标注
- 版本与弃用包含主版本切换方式、弃用标注规则、破坏性变更约束

### 端点目录

- 按资源分组（用户、订单、商品等）
- 表头：方法、路径、认证需求、一句话说明
- 方法使用 HTTP 标准动词（GET、POST、PUT、PATCH、DELETE）

### 端点详情

- 必含"元信息"小节：幂等性、幂等键（如适用）、版本、弃用计划
- 请求体/查询参数/路径参数的字段表格（字段名、类型、必填、约束）
- 成功响应：标注状态码和 JSON 示例
- 错误响应表必含 5 列：状态码、code、**错误类型**（对齐 5 类）、**可重试**、说明

### 错误码参考

- 统一汇总所有错误码，方便全局检索
- 每个错误码必含：HTTP 状态码、业务 code、**错误类型**（5 类之一）、**可重试性**、说明
- 错误类型必须可映射到 `architecture-design` 接口契约规定的 5 类

## 更新规则

| 事件 | 更新内容 |
|------|---------|
| 新增端点 | "端点目录" + "端点详情"（含元信息） |
| 修改端点参数/响应 | 对应端点的"端点详情" |
| 修改端点幂等性 / 版本 / 弃用计划 | 对应端点的"元信息" |
| 新增错误码 | "错误码参考"（含错误类型与可重试性） |
| 认证方式变更 | "概述" + "认证" |
| 通用约定变更 | "通用约定" |
| `architecture-design` 接口契约（错误模型 / 版本策略 / 幂等约定）变更 | "通用约定"对应小节 |

## 与其他文档的关系

- 端点契约具体记录在本文档；**架构层接口决策**（错误模型分类、版本策略、幂等约定、跨进程通信方式）由 `architecture-design` 技能产出
- 本文档「通用约定」必须与 `architecture-design` 接口契约保持一致；冲突时以 `architecture-design` 最新产出为准
- 端点级幂等性、可重试性、错误类型标注必须与本文档「通用约定」与「错误码参考」一致
