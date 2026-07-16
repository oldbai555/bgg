# D2Table 复用收敛 + 大文件拆分（Phase 2 Week 3）

> 前置依赖：`02-domain-reorg-and-api-layer.md` 已完成（本文档按迁移后的新路径描述）。

## 1. D2Table 覆盖现状

47 个视图中 29 个使用 D2Table。以下是明确未使用的文件及初步判断（**执行时必须逐个打开确认，以下判断不是最终结论，只是研究阶段的初筛**）：

| 文件（迁移后路径） | 初筛结论 | 理由 |
|---|---|---|
| `views/iam/DepartmentList.vue` | **合理例外，保留** | 部门是树形数据，D2Table 是扁平分页表格模型，`el-tree` 更合适，与 `MenuList` 同理 |
| `views/iam/MenuList.vue`（409 行） | **合理例外，保留** | 菜单树形结构，同上 |
| `views/monitoring/MonitorList.vue` | **需要在执行时打开核实**，初筛倾向可收敛 | 若数据结构是普通分页列表却手搓表格，应收敛进 D2Table；若有监控特有的实时刷新/图表混排需求，则保留例外并在代码里写明原因 |
| `views/temp/MetricList.vue` | 不适用 | 属于 `views/temp/` 死代码，按 `07` 号文档直接删除，不用讨论是否收敛 |
| `views/public/BlogList.vue`、`views/public/VideoList.vue` | **保留例外** | 公共页面无权限体系、无编辑/删除操作、走独立的"小程序风格"/响应式设计（见 `06`），D2Table 是为管理后台 CRUD 设计的组件，不适合公共展示页，两者本就是不同的产品形态，不应该强行统一 |
| `views/chat/ChatList.vue`（1037 行） | **保留例外，但需要拆分（见下）** | 会话列表是即时通讯 UI（头像+摘要+未读角标+点击进入会话），不是管理后台的 CRUD 表格，套 D2Table 会让交互倒退，例外成立；但文件本身过大，是独立问题（见 §2） |

执行时的判断原则：**D2Table 是为"分页列表 + CRUD 操作"设计的，凡是数据本质上是树形、或页面本质上不是管理后台 CRUD（公共展示页、即时通讯 UI）的，保留例外并在文件顶部写一行注释说明"为什么不用 D2Table"**，不要为了指标好看强行套用。

## 2. `ChatList.vue` 拆分方案（1037 行，全仓最大单文件组件）

执行时先完整读一遍全文件再定具体拆分边界，以下是基于组件命名推测的拆分方向，供起点参考：

1. **列表项渲染**拆成 `components/chat/ChatListItem.vue`（头像、名称、最后一条消息摘要、未读角标、时间戳的展示逻辑）。
2. **搜索/筛选交互**（如果存在）拆成独立的搜索栏子组件或 composable。
3. **WebSocket 事件订阅与本地状态同步逻辑**收进 composable（如 `composables/useChatList.ts`），把"UI 渲染"和"实时数据同步"两个关注点分开，这也是让 `stores/websocket.ts`/`stores/notification.ts`（见 `03` 号文档）拆分后能被 `ChatList.vue` 干净消费的前提。
4. 目标：拆分后主文件（`ChatListPage.vue` 或保留原名）控制在 300 行以内，子组件/composable 各自单一职责。
5. 拆分必须保持行为完全不变（这是纯重构，不是功能改动），拆分前后手工走查一遍会话列表的核心交互（收发消息、未读角标更新、切换会话）确认无回归。

## 3. layout/ 组件调整

`layout/` 下 6 个组件（`AppHeader`/`AppSidebar`/`Breadcrumb`/`MessageNotification`/`PageHeader`/`UserMenu`）本身职责清晰，本轮不做结构性拆分，但需要配合 Phase 3 的暗色模式全面适配（`05` 号文档）和响应式重构（`06` 号文档）做样式层面的改动——具体样式改动不在本文档范围，本文档只确认：**这 6 个组件的组件边界本身是合理的，不需要重新拆分或合并**。

`MessageNotification.vue` 依赖的未读消息数据源，在 `03` 号文档的 store 拆分完成后需要同步改为消费新的 `notification.ts` store（如果拆分后未读状态确实迁移到了独立 store），这是一处需要跨文档协调的改动点，执行 `03` 号文档时顺带检查这个消费点。

## 完成的定义

- 每个未使用 D2Table 的视图文件顶部都有一行注释说明例外原因，或已被收敛进 D2Table。
- `views/chat/ChatList.vue`（或拆分后的主文件）不超过 300 行，拆出的子组件/composable 各自职责单一。
- 拆分前后人工走查会话列表核心交互，行为无回归。
- `npm run typecheck` + `npm run build` 通过。
