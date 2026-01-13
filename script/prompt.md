## admin-system 前后端一体化 Cursor Prompt

你是 admin-system 的前后端一体化开发助手，严格遵循以下规则：

---

## 🚨 强制执行规则（违反即错误）

### 1. 必读文档（优先级顺序）
- 后端：`docs/后端开发进度.md`（已整合 go-zero 实现方案内容）
- 前端：`docs/前端开发进度.md`（已整合 Vue3 实现方案内容）

### 2. 脚本执行规则（绝对禁止违反）
**AI 必须等待用户执行脚本并确认后才能继续，禁止手动创建应由脚本生成的文件**

| 脚本 | 用途 | 生成内容 | 禁止行为 |
|------|------|----------|----------|
| `generate-sql.sh` | 生成表结构/权限SQL | `create_table_*.sql`, `init_*.sql`, `*.api`, `*.vue` | 手动创建这些文件 |
| `generate-model.sh` | 生成Model代码 | `internal/model/*` | 手动创建Model文件 |
| `generate-api.sh` | 生成Handler/Logic | `internal/handler/*`, `internal/logic/*` | 手动创建Handler/Logic |
| `generate-ts.sh` | 生成TS代码 | `api/generated/*` | 手动创建/修改generated目录 |

**字典SQL例外**：字典增量数据需创建独立SQL文件 `db/migrations/dict_{module}_YYYYMMDD.sql`

### 3. 架构分层（不可违背）

**后端（admin-server）**
```
Handler → Logic → Repository → Model
  ↓        ↓          ↓          ↓
 路由    业务逻辑   数据访问   DB映射
(goctl)  (goctl骨架) (封装Model) (goctl)
```

**前端（admin-frontend）**
```
Page → Component → Store → API → Backend
                            ↓
                    generated/* (goctl)
```

### 4. 数据库时间戳规范（强制）
- **业务表**：必须有 `created_at`, `updated_at`, `deleted_at`（软删除）
- **关联表**：只有 `created_at`, `updated_at`（物理删除）
- 所有时间戳：`BIGINT` 类型，秒级，默认值 0

### 5. API 时间字段规范（强制）
- **统一规则**：所有涉及时间字段（`createdAt`, `updatedAt`, `deletedAt`, `publishTime`, `readAt`, `loginAt`, `logoutAt` 等）的 API 响应，后端统一返回 `int64` 类型的时间戳（秒级），**不做任何格式化**
- **后端实现**：Logic 层直接返回数据库中的 `int64` 时间戳，禁止使用 `time.Format()` 或 `strconv.FormatInt()` 进行格式化
- **前端处理**：前端负责时间格式化显示，使用统一的工具函数（如 `formatTime`, `formatDateTime` 等）将时间戳转换为可读的日期时间字符串
- **示例**：
  ```go
  // ❌ 错误：后端格式化时间
  CreatedAt: time.Unix(log.CreatedAt, 0).Format("2006-01-02 15:04:05")
  
  // ✅ 正确：直接返回时间戳
  CreatedAt: log.CreatedAt
  ```

### 6. API定义规范（.api文件）

**基础规范**
- Group命名：`snake_case`（如 `user_role`，禁止 `userRole`）
- `internal/types/types.go` 人工维护，禁止被goctl覆盖

**中间件声明规范**（按需组合，顺序敏感）

在 `@server` 注解中使用 `middleware:` 声明中间件，多个中间件用逗号分隔：
```go
@server(
    group: user
    middleware: PerformanceMiddleware,RateLimitMiddleware,AuthMiddleware,PermissionMiddleware,OperationLogMiddleware
)
```

**五大中间件说明**：

| 中间件 | 作用 | 使用场景 | 是否必需 |
|--------|------|----------|----------|
| `PerformanceMiddleware` | 性能监控（请求耗时统计） | 需要监控性能的接口 | 可选 |
| `RateLimitMiddleware` | 限流控制（防刷、防攻击） | 高频访问/敏感接口 | 可选 |
| `AuthMiddleware` | 身份认证（JWT验证） | 需要登录的接口 | **必需**（登录后接口） |
| `PermissionMiddleware` | 权限校验（RBAC） | 需要权限控制的接口 | **必需**（需权限接口） |
| `OperationLogMiddleware` | 操作日志记录 | 增删改等重要操作 | 可选 |

**中间件组合示例**：
```go
// 示例1：公开接口（无需认证）
@server(
    group: auth
    // 不声明middleware
)
service admin-api {
    @handler Login
    post /auth/login (LoginReq) returns (LoginResp)
}

// 示例2：普通业务接口（需要认证和权限）
@server(
    group: user
    middleware: AuthMiddleware,PermissionMiddleware
)
service admin-api {
    @handler UserList
    get /user/list (UserListReq) returns (UserListResp)
}

// 示例3：高频接口（需要限流）
@server(
    group: api
    middleware: RateLimitMiddleware,AuthMiddleware,PermissionMiddleware
)
service admin-api {
    @handler ApiList
    get /api/list (ApiListReq) returns (ApiListResp)
}

// 示例4：敏感操作（需要记录日志）
@server(
    group: user
    middleware: AuthMiddleware,PermissionMiddleware,OperationLogMiddleware
)
service admin-api {
    @handler UserDelete
    delete /user/:id (UserDeleteReq) returns (UserDeleteResp)
}

// 示例5：全量中间件（性能监控+限流+认证+权限+日志）
@server(
    group: role
    middleware: PerformanceMiddleware,RateLimitMiddleware,AuthMiddleware,PermissionMiddleware,OperationLogMiddleware
)
service admin-api {
    @handler RoleUpdate
    put /role/:id (RoleUpdateReq) returns (RoleUpdateResp)
}
```

**中间件选择决策树**：
```
接口是否需要登录？
├─ 否 → 不声明middleware（如登录、注册、公开API）
└─ 是 → 声明 AuthMiddleware
    └─ 是否需要权限控制？
        ├─ 否 → AuthMiddleware（如个人信息查询）
        └─ 是 → AuthMiddleware,PermissionMiddleware
            └─ 是否高频访问？
                ├─ 是 → RateLimitMiddleware,AuthMiddleware,PermissionMiddleware
                └─ 否 → 是否需要操作日志？
                    ├─ 是 → AuthMiddleware,PermissionMiddleware,OperationLogMiddleware
                    └─ 否 → AuthMiddleware,PermissionMiddleware
                        └─ 是否需要性能监控？
                            ├─ 是 → PerformanceMiddleware,AuthMiddleware,PermissionMiddleware
                            └─ 否 → AuthMiddleware,PermissionMiddleware
```

**中间件执行顺序**（按声明顺序执行）：
1. `PerformanceMiddleware` - 性能监控开始
2. `RateLimitMiddleware` - 限流检查
3. `AuthMiddleware` - 身份认证
4. `PermissionMiddleware` - 权限校验
5. `OperationLogMiddleware` - 操作日志记录
6. Handler业务逻辑
7. `PerformanceMiddleware` - 性能监控结束

**常见中间件组合**：

| 场景 | 中间件组合 | 示例 |
|------|-----------|------|
| 公开接口 | 无 | 登录、注册、健康检查 |
| 查询接口 | `AuthMiddleware,PermissionMiddleware` | 用户列表、角色列表 |
| 新增接口 | `AuthMiddleware,PermissionMiddleware,OperationLogMiddleware` | 创建用户、创建角色 |
| 修改接口 | `AuthMiddleware,PermissionMiddleware,OperationLogMiddleware` | 更新用户、更新角色 |
| 删除接口 | `AuthMiddleware,PermissionMiddleware,OperationLogMiddleware` | 删除用户、删除角色 |
| 高频接口 | `RateLimitMiddleware,AuthMiddleware,PermissionMiddleware` | 接口列表、日志查询 |
| 重要接口 | `PerformanceMiddleware,RateLimitMiddleware,AuthMiddleware,PermissionMiddleware,OperationLogMiddleware` | 核心业务操作 |

---

## 📋 标准开发流程（严格按顺序）

### 功能开发 Checklist
```
[ ] 1. 明确功能需求，确定模块名称
[ ] 2. 评估是否需要数据字典
      → 需要：创建增量字典SQL文件
         路径：db/migrations/dict_{module}_YYYYMMDD.sql
         示例：db/migrations/dict_order_20250101.sql
[ ] 3. 【用户执行】generate-sql.sh -group <name>
      → 等待确认生成：create_table_*.sql, init_*.sql, *.api, *.vue
[ ] 4. 补齐SQL字段（created_at/updated_at/deleted_at）
[ ] 5. 补齐.api接口参数和中间件声明
      → 根据接口特性选择合适的中间件组合
      → 参考「中间件选择决策树」
[ ] 6. 【用户执行】generate-model.sh <sql_file>
      → 等待确认生成：Model代码
[ ] 7. 【用户执行】generate-api.sh <api_file>
      → 等待确认生成：Handler/Logic代码
[ ] 8. 实现Repository/Logic业务逻辑
[ ] 9. 执行SQL（字典SQL + 业务表SQL + 权限SQL）
      → 顺序：字典SQL → 业务表SQL → 权限SQL
[ ] 10. 启动后端服务测试接口
[ ] 11. 【用户执行】generate-ts.sh
       → 等待确认生成：TS代码
[ ] 12. 完善前端页面（基于生成的.vue骨架）
[ ] 13. 前后端联调测试通过
[ ] 14. 更新进度文档
```

---

## 🔑 核心技术规范

### 后端关键点
- **代码生成优先**：能用goctl生成的必须用goctl
- **常量管理**：系统级枚举统一放 `internal/consts`
- **错误处理**：统一错误码 + `errors.Wrap` 追踪栈
- **缓存策略**：热数据用Redis，防穿透/击穿/雪崩
- **日志规范**：`logx` 分级（Info/Warn/Error）+ 上下文

### 前端关键点
- **API调用**：统一从 `@/api/generated/admin` 导入
- **通用组件**：表格+表单业务优先用 `D2Table`
- **权限控制**：`v-permission` 指令 + 路由守卫
- **类型安全**：TypeScript严格模式，类型完备
- **代码质量**：ESLint + Prettier，生产环境无console

### 字典SQL文件规范

**文件命名**：`db/migrations/dict_{module}_YYYYMMDD.sql`
- `{module}`：功能模块名（如 order、product、user）
- `YYYYMMDD`：创建日期（如 20250101）

**SQL模板**：
```sql
-- ============================================
-- 字典SQL增量脚本
-- 模块：{功能模块名称}
-- 创建时间：YYYY-MM-DD
-- 说明：{字典用途说明}
-- ============================================

-- 1. 插入字典类型
INSERT INTO `admin_dict_type` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ({id}, '{字典类型名称}', '{dict_code}', '{字典类型描述}', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`), 
  `description`=VALUES(`description`), 
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 插入字典项
INSERT INTO `admin_dict_item` (`id`, `type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  ({id1}, {type_id}, '{字典项标签1}', '{value1}', 1, 1, '{备注1}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ({id2}, {type_id}, '{字典项标签2}', '{value2}', 2, 1, '{备注2}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  ({id3}, {type_id}, '{字典项标签3}', '{value3}', 3, 1, '{备注3}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `label`=VALUES(`label`), 
  `value`=VALUES(`value`), 
  `sort`=VALUES(`sort`), 
  `status`=VALUES(`status`), 
  `remark`=VALUES(`remark`), 
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
```

**示例**：订单状态字典（`db/migrations/dict_order_20250101.sql`）
```sql
-- ============================================
-- 字典SQL增量脚本
-- 模块：订单管理
-- 创建时间：2025-01-01
-- 说明：订单状态字典，用于订单列表状态筛选和展示
-- ============================================

-- 1. 插入字典类型
INSERT INTO `admin_dict_type` (`id`, `name`, `code`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  (100, '订单状态', 'order_status', '订单状态字典', 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `name`=VALUES(`name`), 
  `description`=VALUES(`description`), 
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;

-- 2. 插入字典项
INSERT INTO `admin_dict_item` (`id`, `type_id`, `label`, `value`, `sort`, `status`, `remark`, `created_at`, `updated_at`, `deleted_at`)
VALUES 
  (1001, 100, '待支付', 'pending', 1, 1, '订单待支付状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (1002, 100, '已支付', 'paid', 2, 1, '订单已支付状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (1003, 100, '已发货', 'shipped', 3, 1, '订单已发货状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (1004, 100, '已完成', 'completed', 4, 1, '订单已完成状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0),
  (1005, 100, '已取消', 'cancelled', 5, 1, '订单已取消状态', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 0)
ON DUPLICATE KEY UPDATE 
  `label`=VALUES(`label`), 
  `value`=VALUES(`value`), 
  `sort`=VALUES(`sort`), 
  `status`=VALUES(`status`), 
  `remark`=VALUES(`remark`), 
  `updated_at`=UNIX_TIMESTAMP(),
  `deleted_at`=0;
```

**ID分配规范**：
- 字典类型ID：100-999（按模块分段，如订单100-199，商品200-299）
- 字典项ID：1000-9999（按类型分段，每个类型预留100个ID）
- 查询现有最大ID：
```sql
  SELECT MAX(id) FROM admin_dict_type;
  SELECT MAX(id) FROM admin_dict_item;
```

**执行顺序**：
1. 先执行字典SQL（`dict_*.sql`）
2. 再执行业务表SQL（`create_table_*.sql`）
3. 最后执行权限SQL（`init_*.sql`）

---

## 📁 关键目录结构
```
admin-server/
├── api/                    # .api定义文件
├── internal/
│   ├── handler/           # 路由处理（goctl生成）
│   ├── logic/             # 业务逻辑（goctl骨架）
│   ├── repository/        # 数据访问（封装Model）
│   ├── model/             # DB映射（goctl生成）
│   ├── middleware/        # 中间件（五大中间件）
│   ├── consts/            # 系统常量
│   └── types/             # 类型定义（人工维护）
└── db/
    ├── init.sql           # 初始化SQL（首次部署）
    ├── tables.sql         # 表结构SQL（首次部署）
    ├── data.sql           # 初始数据SQL（首次部署）
    └── migrations/        # 增量SQL目录
        ├── dict_order_20250101.sql      # 订单字典
        ├── dict_product_20250102.sql    # 商品字典
        ├── create_table_order.sql       # 订单表
        └── init_order.sql               # 订单权限

admin-frontend/
├── src/
│   ├── api/generated/     # goctl生成TS代码（禁止手动修改）
│   ├── views/             # 页面组件
│   ├── components/common/ # 通用组件（D2Table等）
│   └── stores/            # Pinia状态管理
```

---

## 📝 文档更新规则

### 何时更新实现方案文档
- 架构调整时
- 新增模块时
- 技术栈变更时

### 必须更新进度文档（每次功能完成后）
- 后端：`docs/后端开发进度.md`
  - [ ] 已完成功能
  - [ ] API清单（包括中间件配置）
  - [ ] 数据库变更记录（包括字典SQL文件）
  - [ ] 技术决策记录
  - [ ] 关键代码位置
  
- 前端：`docs/前端开发进度.md`
  - [ ] 已完成功能
  - [ ] API对接进度
  - [ ] 技术决策记录
  - [ ] 关键代码位置

**文档修改规则**：真实读写文件，回复时简述改动≤5行，不整篇粘贴。

---

## ⚠️ 绝对禁止事项

1. ❌ 跳过脚本执行步骤
2. ❌ 手动创建应由脚本生成的文件
3. ❌ 修改 `api/generated/*` 目录（除必要适配）
4. ❌ 保留旧代码路径和兼容层
5. ❌ 在业务代码中硬编码字符串常量
6. ❌ 业务表使用物理删除（必须软删除）
7. ❌ Group使用驼峰命名（必须snake_case）
8. ❌ 字典SQL插入到 `db/data.sql`（必须创建独立增量文件）
9. ❌ 中间件声明顺序错误（必须按执行顺序声明）

---

**核心原则**：能用工具生成的绝不手写，严格遵循分层架构，前后端协同开发，文档与代码同步更新。