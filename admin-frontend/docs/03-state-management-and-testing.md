# 状态管理拆分 + vitest 测试基建（Phase 1 Week 2 收尾 + Phase 2 Week 4）

> 前置依赖：`01-architecture-target.md`（类型安全部分，测试需要类型收敛后的 `Envelope<T>` 才好写）。vitest 基建部分属于 Phase 1 Week 2，store 拆分与组件层测试补充属于 Phase 2 Week 4，同一份文档覆盖两个时间点是因为两者都是"状态管理与测试"这一个技术主题，拆成两篇反而增加交叉引用成本。

## 1. Store 审计结论

| Store | 行数 | 结论 |
|---|---|---|
| `stores/app.ts` | 55 | 单一职责（主题/折叠状态等全局 UI 配置），不动 |
| `stores/dict.ts` | 123 | 单一职责（字典缓存 + `REQUIRED_DICT_CODES`），不动 |
| `stores/user.ts` | 113 | 单一职责（登录态/用户信息/菜单/权限），不动 |
| `stores/websocket.ts` | 428 | **拆分**：现状同时持有连接生命周期（connect/reconnect/心跳）+ 未读消息列表 + `unreadCount` 等聚合 getters，是本轮唯一的"胖 store"问题 |

### `websocket.ts` 拆分方案

先在 Phase 1 Week 2 做类型/wrapper 相关的收尾工作，`websocket.ts` 的实际拆分放到 Phase 2 Week 4（此时 Phase 1 的域重组已完成，chat 域目录已经就位，拆分出的"未读消息"状态可以合理地与 `views/chat/` 域邻近，减少来回改动）：

1. **连接生命周期职责**保留在 `stores/websocket.ts`：`connect`/`disconnect`/重连策略/心跳/`readyState`，这是真正意义上跨应用生命周期的全局连接状态，适合留在 Pinia store。
2. **未读消息状态**（未读列表、`unreadCount` 等派生 getters）拆到新的 `stores/notification.ts`（若与现有 `system/NotificationList.vue` 使用的通知概念本就相关，合并成一个 store 更自然，避免"未读消息"和"系统通知未读"两套并行的未读计数逻辑）——**具体是否合并，执行时先读一遍 `websocket.ts` 全文和 `NotificationList.vue`/`MessageNotification.vue` 的实际字段命名再定，不要凭本文档的描述臆断字段名**。
3. `websocket.ts` 对外暴露的连接实例通过 `provide/inject` 或直接 import store 的方式让新 store 订阅消息事件更新未读状态，两个 store 之间是"连接 store 广播事件，通知 store 消费"的单向依赖，不要反向依赖。
4. 拆分后每个 store 控制在 150-200 行以内，若发现拆完仍然偏大，说明还有职责没分干净，不要为了凑数字而强行拆分。

## 2. vitest 测试基建（Phase 1 Week 2）

### 配置改动点

- `package.json` 新增 `devDependencies`：`vitest`、`@vue/test-utils`、`jsdom`（或 `happy-dom`，二选一，`happy-dom` 更轻量，除非遇到兼容性问题否则优先用它）、`@vitest/coverage-v8`（可选，仅用于本地查看覆盖率，不接入 CI 门槛）。
- `package.json` `scripts` 新增 `"test": "vitest run"`、`"test:watch": "vitest"`。
- `vite.config.ts` 新增 `test` 配置块（`environment: 'happy-dom'`、`globals: true` 可选、`setupFiles` 如需要）；vitest 复用 Vite 的别名解析（`@/` 指向 `src/`），不需要额外配置路径映射。
- 不新增独立的 `vitest.config.ts`，直接在现有 `vite.config.ts` 里扩展 `test` 字段，减少一份配置文件。

### 测试覆盖范围优先级（不追求覆盖率门槛，按优先级做到"值得测的都测了"为止）

1. **stores**：`dict.ts`（字典加载/缓存/TTL 过期逻辑）、`user.ts`（token 存取、权限判断）、拆分后的 `websocket.ts`/`notification.ts`（连接状态流转、未读计数逻辑）——这些是纯逻辑 + 少量浏览器 API（`localStorage`），最适合单测，且历史上这类逻辑改一次容易忘了改另一处。
2. **composables**：`useDictOptions.ts`、`useAppConfig.ts`、`usePermission.ts`（合并到 `composables/` 后）——纯函数/纯逻辑，无 UI 依赖。
3. **`utils/request.ts` 拦截器**：`01` 号文档改造后的 `isEnvelope` 类型守卫 + 成功/失败/`10003` 过期码分支逻辑，用 mock axios 实例验证行为不因类型改造而变化（这组测试也是 `01` 号类型安全改造的回归保护网，建议先写测试再动 `request.ts`，即改造前先固化当前行为的测试用例）。
4. **纯函数 utils**：`src/utils/` 下其余不依赖组件生命周期的工具函数（如时间格式化、`generateBreadcrumb` 等）。
5. **关键组件**（按需，不强求）：`D2Table.vue` 的分页/事件触发逻辑用 `@vue/test-utils` 做浅层测试；`router/index.ts` 里 `01` 号文档 A.5 新抽取的 `generateUniqueRouteName` 纯函数应该有测试，因为它是"消灭重复逻辑"改造的直接产物，值得一个测试防止回归。

### 不做的事

- 不引入 Cypress/Playwright 之类的浏览器端 E2E 框架（`00` 号文档 §6 已明确），浏览器走查用 `webapp-testing` skill 人工/半自动验证。
- 不对 47 个视图逐一写组件快照测试——视图层大多是 D2Table 配置 + 简单表单，逻辑集中在 API 调用和字典渲染，真正值得测的逻辑已经下沉到 stores/composables，视图层测试投入产出比低。
- 不设 CI 覆盖率百分比门槛（`00` 号文档已明确）。

## 完成的定义

- `npm run test` 可执行且全绿。
- `stores/`、`composables/` 下所有文件都有对应的 `*.spec.ts`（同目录或 `__tests__/` 子目录均可，与 `@vue/test-utils` 官方推荐一致，选一种全仓统一）。
- `websocket.ts` 拆分后，原文件与新 `notification.ts`（或最终确定的名称）均不超过 200 行，职责边界在文件头部用一行注释说明（不写多段落 docstring，遵循项目"默认不写注释，非显而易见的 WHY 才写"的原则——这里注释是为了防止未来又把两个职责揉回一起）。
