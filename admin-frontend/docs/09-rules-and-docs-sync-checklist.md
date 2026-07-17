# 规则/文档同步清单（贯穿 Phase 1-3，Phase 3 Week 7 做最终收尾）

> 用户明确要求："别忘了对应的文档的修改也需要纳入修改计划。例如 `.cursor/rules`、`AGENTS.md`。"本文档对齐 `admin-server/docs/13-rules-sync-checklist.md` 的先例：本轮改动会让现有规则/文档的哪些具体段落过时，逐条列出，不是"最后随便扫一眼"。**同步不是收尾才做的事**——每个 Phase 结束都要回来过一遍这张表里对应 Phase 的条目，Phase 3 收尾做一次最终的实质性重写确认全部一致。

## 1. 需要同步的文件清单与触发条件

| 文件 | 触发它需要改的原因 | 何时改 |
|---|---|---|
| `.claude/rules/20-frontend.md` | 目录结构一节（"技术栈与目录结构"）描述的 `views/` 按 `system/blog/sdk/video/chatroom/public` 分域，重组后变成 `iam/system/monitoring/misc/content/chat/sdk/task/public`；"API 层规范"一节需要加入"禁止直接 import generated 里的函数，一律走 wrapper，`import type` 例外"的强制规则（目前规则只说"业务代码统一从二次封装层导入"但没有强制且现状违反，本轮要把这条从"建议"变成有落地事实支撑的"强制"）——**"禁止直接 import 函数"与"`import type` 例外"必须作为同一批次一起写入，不能只写强制规则、把类型例外拖到后面，否则中间态会让执行者误以为类型 import 也被禁止**；组件与状态管理一节需要补充 composables/hooks 合并后的规则（`hooks/` 目录不再存在） | Phase 1 结束后（域重组+API层完成时）先改一版，Phase 3 视觉部分完成后再补一版 |
| `.claude/rules/21-public-pages.md` | 现有"小程序风格"DOM/样式契约（`public-list-page`/`.list-grid`/768px 断点等）会被 `06` 号文档的响应式重构整体替换 | Phase 3 公共页面重构完成后（不要在重构前改，避免文档和代码不同步） |
| `.cursor/rules/20-frontend.mdc`、`.cursor/rules/21-public-pages.mdc` | `00-workflow.md` 明确"两者内容保持同步，更新规则时两处都要改"——`.cursor/rules/*.mdc` 是 SSOT，改动从这里发起，再同步到 `.claude/rules/*` 对应文件 | 与上面两条同批次 |
| `AGENTS.md` | 是"面向任意 AI 工具的整合版操作手册"，包含前端目录结构、API 层规范等与 `.claude/rules`/`.cursor/rules` 重复的内容，两处的更新必须同步，不能只改规则文件忘了 `AGENTS.md`（`00-workflow.md` 本身也是 `AGENTS.md` 引用的规则文件之一，是同一份规则体系的不同呈现层） | 与规则文件同批次修改，逐条核对不要遗漏 |
| `docs/前端开发进度.md`（2026-07-17 起已退场归档为 `docs/changelog/archive-frontend.md`，本行为撰写时点的真实记录，不回溯改写） | §0"目录结构要点"、§1"核心功能索引"、§5"关键代码位置"三节列出的具体文件路径（如 `src/views/system/{RoleList,...}.vue`）会因域重组全部失效；§2"已完成功能"提到的"demo 功能页面：`DemoList.vue`"如果按 `07` 号文档删除，对应描述要更新为历史记录（不是删掉这段历史，而是标注"已废弃"，因为这是过程记录不是当前状态清单——具体怎么处理参考该文档"这是历史/参考日志"的定位）；`docs/前端开发进度.md` §4"技术决策记录"应追加本轮重构的关键决策条目（不是搬 `admin-frontend/docs/progress.md` 的全部内容，只记功能行为发生实质变化的部分，遵循 `00-refactor-overview.md` 里已经明确的两份文档分工原则） | 各 Phase 结束时增量更新对应部分，不要攒到最后一次性改 |
| `admin-server/scripts/sqlgen/templates/list_page.vue.tpl` | **已核实的具体冲突**：模板第 47 行 `import { {{.GroupFuncName}}List, ... } from '@/api/generated/admin'` 直接从 generated 导入，与 `01`/`02` 号文档"视图禁止直接 import generated 函数，必须走域 wrapper"的新规则直接冲突——新脚手架生成的页面从第一天起就违反新规则 | Phase 1 域重组 + API wrapper 规范落地后**立即**改（不要拖到 Phase 3），否则 Phase 1-3 期间新增的任何标准 CRUD 模块都会继续生成违规代码，扩大后续清理面 |
| `admin-server/scripts/sqlgen/templates/init_module.api.tpl` | 需要核实：脚手架生成的 `.api` 草稿默认 `group` 命名是否已经隐含假设了旧的前端目录分组（若无强绑定则不用改，若草稿里有前端路径假设需要一并核实） | Phase 1，与上一条同批次核实 |

## 2. 模板改法（`list_page.vue.tpl` 具体改动点）

按 `00-workflow.md`"发现生成骨架的默认写法不符合项目最新约定，应该改模板而不是每次生成后手工修一遍"的既有原则，改法：

1. 第 47 行的 import 改为从 `@/api/{{.Domain}}`（新增模板变量，对应 `02` 号文档的 8 个域 wrapper 之一）导入，而不是硬编码 `@/api/generated/admin`；模板变量 `.Group` 目前是 `<domain>/<module>` 格式（如 `iam/user`），可以直接从中解析出 `.Domain` 前缀部分复用，不需要新增用户输入。
2. 若域 wrapper 尚未包含该次新生成模块对应的函数（例如新增了一个全新域的模块），模板生成后需要执行者手动在对应 `api/<domain>.ts` 里补一个函数出口——这是脚手架固有的"生成骨架，人工补业务细节"模式的延伸，不是缺陷。
3. 类型 import（`{{.GroupUpper}}Item` 等）保留从 `@/api/generated/admin` 直接导入，符合 `01` 号文档 A.2 的类型 import 例外规则，不用改。

## 3. Cursor 端同步动作

`00-workflow.md` 已有的维护流程：改完 `.cursor/rules/*.mdc`（SSOT）后执行 `make sync-claude-rules` 同步生成 `.claude/rules/*`。本轮每次改动规则文件后都要跑这条命令，不要让两端漂移——这条动作本身在"可以直接执行"范围内（`08` 号文档 §1），不需要为此单独确认。

## 完成的定义（Phase 3 Week 7 最终收尾时核对）

- `.claude/rules/20-frontend.md`、`.claude/rules/21-public-pages.md`、对应的 `.cursor/rules/*.mdc`、`AGENTS.md` 四处关于前端目录结构/API 规范/公共页面契约的描述与重构后的实际代码一致，交叉核对无遗漏字段。
- `docs/前端开发进度.md`（现已归档为 `docs/changelog/archive-frontend.md`）的目录结构、功能索引、关键代码位置三节路径已在归档前更新为重构后的新路径。
- `list_page.vue.tpl`（及核实后视情况处理的 `init_module.api.tpl`）已改为符合新 API wrapper 规范，用一次新增模块的脚手架生成验证产出代码不再直接 import generated。
- `make sync-claude-rules` 已执行，两端规则文件无漂移。
