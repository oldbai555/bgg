# Part A：目标架构技术决策（Phase 1 总纲）

> 前置依赖：读完 `00-refactor-overview.md`。本文档是 Phase 1（架构地基）的技术总纲，`02`（域重组+API层）、`03`（状态管理+测试）的前半部分、`07`（工具链清理）都是它的可执行拆解。遇到需要复核设计决策的地方回指本文档，不要在分篇文档里重新推导一遍。

## A.1 域目录规范

### 目标目录结构

```
src/
├── api/
│   ├── generated/          # goctl 生成，禁止手改（不变）
│   ├── iam.ts              # 新增：用户/角色/权限/部门/菜单/API 管理
│   ├── system.ts           # 新增：配置/字典/文件/公告/通知
│   ├── monitoring.ts       # 新增：审计日志/登录日志/性能日志/监控 + 现有 metric.ts 合并进来
│   ├── misc.ts              # 新增：demo/daily_short_sentence（若这两个功能仍保留，见 07 号文档死代码清单）
│   ├── content.ts           # 新增：合并现有 blog.ts + video.ts
│   ├── chat.ts               # 新增：会话/消息/群组
│   ├── sdk.ts                 # 新增：API Key/接口管理/调用日志
│   ├── task.ts                 # 新增：异步任务
│   └── public.ts               # 保留：字典等技术性/跨域公共接口
├── views/
│   ├── iam/                # UserList, RoleList, PermissionList, DepartmentList, MenuList, ApiList
│   ├── system/              # ConfigList, DictTypeList, DictItemList, FileList, NoticeList, NotificationList
│   ├── monitoring/           # AuditLogList, LoginLogList, MonitorList, OperationLogList, PerformanceLogList, MetricStats
│   ├── content/                # 原 blog/ + video/ 全部文件
│   ├── chat/                    # 原 chatroom/ 改名（与后端 chat-rpc 命名对齐，去掉不必要的 "room" 后缀）
│   ├── sdk/                      # 不变
│   ├── task/                      # 新增，TaskList.vue 从 system/ 移入
│   └── public/                     # 不变（域内不再细分，公共页在 06 号文档单独重构）
```

### 命名与归属规则

- 目录名与 `api/<domain>.ts` 一一对应，且与后端 `admin.api` 的 `group:` 前缀（`iam/user`、`system/config` 等）保持同一套域名词汇——工程师能直接从后端 group 名猜到前端文件应该放哪，不需要再查映射表。
- `misc` 域是否单独建目录取决于 `demo`/`daily_short_sentence` 两个功能的去留结论（见 `07-cleanup-and-tooling.md`）；如果确认删除演示性质的 `DemoList`，`misc` 域可以只留 `api/misc.ts`（如果还有其他 misc 接口需要调用）或直接不建。
- `views/chatroom/` 改名为 `views/chat/` 属于目录重命名，不是新建业务；迁移时用 `git mv` 保留历史。
- 具体到每个文件的迁移目标见 `02-domain-reorg-and-api-layer.md`，本文档只定规则。

## A.2 API 层规范

### 现状问题

`grep "from '@/api/generated"` 命中 85 处，其中相当一部分是函数调用（`iamApi.userList(...)` 这种直接调 generated 的模式），违反 `.claude/rules/20-frontend.md` "业务代码只从二次封装层导入"的既有规则。只有 blog/video/metric/public 四个域有 wrapper，其余五个域（iam/system/monitoring/chat/sdk/task）完全没有。

### 目标规则

1. **每个域一个 wrapper 文件**（`src/api/<domain>.ts`），职责与现有 `blog.ts` 一致：包一层调用 `generated/` 里的请求函数，做错误处理/拦截器集成/统一返回类型，修正生成路径里多余的前缀问题。
2. **视图/组件禁止直接 `import` `generated/` 里的请求函数**（如 `iamApi.userList`），一律改为 `import { userList } from '@/api/iam'`。
3. **类型 import 例外**：`import type { UserItem } from '@/api/generated/admin'` 允许直接从 generated 引入，不强制在 wrapper 里重新导出——这是本轮和 admin-server 后端 wrapper 规则的一个务实差异（后端 RPC client 类型收敛成本低，前端为每个类型建 re-export 纯增加维护面，无实质收益）。
4. wrapper 函数命名与 generated 保持一致（不重新发明命名），只是加上域前缀的模块归属，调用处从 `import { xxxApi } from '@/api/generated/admin'; xxxApi.yyy()` 改成 `import { yyy } from '@/api/xxx'`。

### 迁移方式

`02-domain-reorg-and-api-layer.md` 给出每个域 wrapper 里需要包含哪些函数、以及现有 85 处直接引用 generated 的调用点如何逐一改写为走 wrapper。

## A.3 类型安全：Envelope\<T\> 收敛

### 现状问题（`src/utils/request.ts`）

- 第 39/40/43/67 行共 4 处 `as any`：`typeof (res as any).code === 'number'`、`const code = (res as any).code`、`return (res as any).data`、`(res as any).msg`。
- 拦截器 `resp` 回调返回类型实际上是 `any`（因为 axios 拦截器签名要求返回 `resp` 或 `resp.data`，这里返回的是解包后的 `data`，类型完全丢失），导致所有业务代码里 `await xxxApi.yyy()` 拿到的返回值类型都是隐式 `any`，`.claude/rules/20-frontend.md` 要求的类型安全在请求这一层就已经破防。

### 目标方案

1. 定义显式的 `Envelope<T>` 类型（新建 `src/types/envelope.ts` 或加入现有 `src/types/`）：
   ```ts
   interface Envelope<T> {
     code: number
     msg: string
     data: T
   }
   ```
2. 给 axios 实例的响应拦截器加类型：用 `AxiosResponse<Envelope<unknown>>` 而不是隐式 `any`，用类型守卫函数 `isEnvelope(res: unknown): res is Envelope<unknown>` 替代当前的 `'code' in res && typeof (res as any).code === 'number'` 判断，消灭全部 4 处 `as any`。
3. 保持现有"非标准结构直接返回原始 data（兼容字典等特殊接口）"的行为不变（第 70-71 行注释说明的场景），只是用类型守卫代替裸转换，不改变运行时行为——这是纯类型安全改造，不是行为改造。
4. axios 拦截器本身的返回类型受限于 axios 的类型定义（`InternalAxiosResponse` 的返回类型语义特殊），无法做到完全端到端的编译期类型推导（这是 axios 拦截器模式的已知限制，不是本次改造能力所限）；改造目标是**消灭现有的 4 处显式 `as any`**，不是重写整个 HTTP 客户端为 fetch 或引入类型更强的库（如 `ky`/`openapi-fetch`）——那是过度设计，超出本轮范围。

## A.4 composables/ 与 hooks/ 合并

### 现状

`composables/` 下 `useAppConfig.ts`、`useDictOptions.ts`；`hooks/` 下仅 `usePermission.ts`。三个文件体量都不大，没有文档化的语义分工（不是"业务逻辑 vs UI 逻辑"这种刻意区分），是历史遗留的偶然分裂。

### 目标方案

- 统一合并到 `composables/`（Vue 3 生态更通用的叫法），删除 `hooks/` 目录。
- `usePermission.ts` 移动到 `composables/usePermission.ts`，所有 `import {usePermission} from '@/hooks/usePermission'` 改为 `from '@/composables/usePermission'`（`router/index.ts` 第 3 行等处需要同步改）。
- 不引入新的目录拆分规则（如再分 `composables/business/` vs `composables/ui/`）——三个文件的规模不值得再分层，避免过度设计。

## A.5 路由组件映射修复

### 现状问题（`src/router/index.ts`）

- 第 133-152 行 `resolveComponent(component, path)`：对 `import.meta.glob('../views/**/*.vue')` 产出的 `viewModules` 做字符串清洗（去掉 `../views/` 前缀和 `.vue` 后缀）后与后端菜单下发的 `component`/`path` 字段做数组 `includes` 比对，没有编译期校验；找不到时只 `console.error`（第 149 行）并让该菜单项静默不可达，用户点击后无任何反馈。
- 第 154-227 行 `buildRoutesFromMenus`：第 164-176 行（页面类型菜单）与第 196-208 行（目录类型菜单）几乎逐行重复"生成唯一 routeName"的 while 循环逻辑。

### 目标方案

1. **抽取公共函数** `generateUniqueRouteName(rawPath: string, usedNames: Set<string>): string`，消灭 `buildRoutesFromMenus` 内部两处重复的 while 循环。
2. **组件映射从"运行时字符串比对"改为"域重组后天然对齐的路径约定 + 启动期校验"**：
   - 域目录重组后（见 A.1），前端目录结构（`views/<domain>/<Page>.vue`）应该与后端菜单表里存储的 `component` 字段值形成**约定优于配置**的一一对应关系，不需要额外维护映射表。
   - `resolveComponent` 保留基于 `import.meta.glob` 的动态导入机制（Vite 的能力，继续沿用，不是本次改造重点），但要在应用启动时（`main.ts` 或路由初始化阶段，仅 dev 环境）主动比对"当前已知的后端菜单 `component` 取值集合"与"`viewModules` 的 key 集合"，把不匹配项作为**启动期警告**而不是运行时静默 `console.error`——具体做法：登录后 `fetchMenus()` 拿到菜单数据时，遍历一次做校验并用 `ElMessage.warning`（仅 dev 环境）或至少更醒目的日志提示，而不是留给用户点击菜单才发现 404。
   - 不引入构建期代码生成（如根据后端菜单表生成路由清单的脚本）——这超出本轮范围，且后端菜单是运行时数据（用户可在菜单管理页面新增），不是编译期可知的静态集合，构建期强校验本质上做不到，只能做运行时的"尽早暴露"。
3. 该修复必须与域目录重组（A.1）同批完成，避免目录改名后菜单 `component` 字段短暂失配导致的中间态混乱——具体执行顺序见 `02-domain-reorg-and-api-layer.md` 的"路由 meta / 菜单 `component` 字段联动迁移步骤"。

## 完成的定义

- `npm run typecheck` 通过，`src/utils/request.ts` 中 `as any` 数量归零。
- `composables/`、`hooks/` 合并完成，`hooks/` 目录不存在，全仓无 `from '@/hooks/'` 残留 import。
- `router/index.ts` 的 `buildRoutesFromMenus` 内部无重复的 routeName 生成逻辑。
- 以上改动仅涉及 `request.ts`、`router/index.ts`、`composables/`/`hooks/` 目录、新建的 `src/types/envelope.ts`，不涉及域目录重组本身（那是 `02` 号文档的范围）。
