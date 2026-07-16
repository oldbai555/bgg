# 13 · `.cursor/rules` ↔ `AGENTS.md` 同步清单

## 前置依赖

无独立前置，但内容会随 Phase 1/3 推进逐步生效——本篇现在只做"找出哪些段落会过期 + 分批规划"，具体改字这个动作要等对应 Phase 真正落地后再执行，不要提前改。

## 0. 已核查的当前状态

`AGENTS.md` 与 `.cursor/rules/00-workflow.mdc`/`.cursor/rules/10-go-code-style.mdc` 内容已通读比对，二者当前基本同源（`CLAUDE.md` 要求的"两处保持同步"目前是满足的）。以下列出的每一条都注明了**具体文件 + 具体行号**，不是泛泛而谈。

## 1. Phase-1-end 同步批次（scripts 相关，范围小，Phase 1 收尾时做）

### 1.1 `AGENTS.md` 第 28 行 / `.cursor/rules/10-go-code-style.mdc` 第 23-25 行：`internal/domain/{iam,task}` 会扩展成 5 个分组

**现状**：
- `AGENTS.md` 第 28 行："`internal/domain/{iam,task}/`（领域服务）"
- `AGENTS.md` 第 71 行："复杂横切逻辑仅在 `internal/domain/iam/`（RBAC）、`internal/domain/task/`（调度）两处引入领域服务"
- `.cursor/rules/10-go-code-style.mdc` 第 23-25 行：
  ```
  ├── domain/
  │   ├── iam/permission_resolver.go  # RBAC 领域服务
  │   └── task/                       # 任务调度领域服务
  ```

**为什么会过期**：A.2 明确要求 `internal/domain/<domain>` 从 Phase 1 一开始就按 Part B 最终的 5 个服务边界分组（`iam` 含 monitoring/system/misc、`content` 含 blog+video、`chat`、`task`、`sdk` 各自独立），不再是"只有 iam、task 两个域有领域服务"。Phase 1 Week 2-3 落地 IAM/Chat/Task 之外的域（blog/video/sdk、monitoring/system/misc）之后，这两处"仅 iam、task 两处"的表述就不准确了——届时会出现 `internal/domain/chat/`、`internal/domain/sdk/`、`internal/domain/content/`（或等价分组名，具体命名以 `01-architecture-target.md`/`04`-`07` 号文档落地时确定的为准）等新目录，且各分组内部只包含"真正需要领域服务"的约 35-40 个方法，不是所有 9 个域都 1:1 对应一个 `internal/domain/<x>/` 目录。

**Phase 1 结束时要改成什么**：把这三处"仅 iam、task 两处"的措辞，改为"按 5 个未来服务分组（iam/content/chat/task/sdk）组织，只对满足分层标准的方法建领域服务"，并列出 Phase 1 结束时实际存在的 `internal/domain/` 子目录清单（具体清单等 Phase 1 Week 2-3 执行完之后填，本篇不预先猜测最终目录名）。

### 1.2 `AGENTS.md` 第 72 行 / `.cursor/rules/10-go-code-style.mdc` 第 46 行："旧代码可内联...直至按域迁移"这句话会失效

**现状**：
- `AGENTS.md` 第 72 行："Logic 优先 `l.svcCtx.Domain.IAM.User`，旧代码可内联 `xxxrepo.NewXxxRepository(svcCtx.Repository)` 直至按域迁移"
- `.cursor/rules/10-go-code-style.mdc` 第 46 行："Logic 优先 `l.svcCtx.Domain.IAM.User`，旧代码可继续内联 `xxxrepo.NewXxxRepository(svcCtx.Repository)` 直至按域迁移"

**为什么会过期**：这句话描述的是"迁移期间的过渡容忍"——当前 161 个 logic 文件里只有 1 个真正用上 `svcCtx.Domain`，其余仍在内联构造 repository，这句话是在给这种过渡状态开绿灯。但本轮重构的目标就是把这个过渡状态**结束掉**：Phase 1 Week 1-3 会把全部 logic 文件迁移到统一走 `svcCtx.Domain.X.Y`（简单域）或领域服务方法（复杂域），迁移完成后"旧代码可内联"这个例外条款就不再有效——如果 Phase 1 结束后这句话还留着，会变成给后续新代码开的一个不该存在的后门（新写的 logic 文件本不该再有理由直接内联构造 repository）。

**Phase 1 结束时要改成什么**：删除"旧代码可内联...直至按域迁移"这个分句，改为硬性要求："Logic 一律通过 `svcCtx.Domain.<Group>.<X>` 访问 repository / 领域服务，禁止在 logic 方法体内内联 `xxxrepo.NewXxxRepository(...)`"。这一条要同时在 `AGENTS.md` 第 72 行和 `.cursor/rules/10-go-code-style.mdc` 第 46 行改，两处措辞保持一致。

### 1.3 `.cursor/rules/10-go-code-style.mdc` 第 88 行：squirrel 迁移"已知例外"实际上已经是过期路径

**现状**：第 88 行原文："**已知例外（待办，不是可以效仿的先例）**：`internal/repository/performance_log_repository.go`、`internal/repository/chat_repository.go` 目前仍未迁移到 squirrel。"

**核查结果（额外发现，不在计划原文里，但已实地验证）**：这两个文件当前的真实路径分别是 `internal/repository/monitoring/performance_log_repository.go`、`internal/repository/chat/chat_repository.go`（已经在上一轮 DDD-lite 重构里挪进域子目录了），第 88 行写的扁平路径 `internal/repository/performance_log_repository.go`/`internal/repository/chat_repository.go` **本身就已经是错的**，与本轮重构无关，是历史遗留的文档滞后。另外用 `grep -c "fmt.Sprintf"` 核查这两个文件，均为 0 次匹配（没有 `fmt.Sprintf` 拼 SQL），但也确认没有使用 `sq.`/squirrel——具体用的是什么方式构建 SQL 需要打开文件确认（不在本篇核查范围内），第 88 行"仍未迁移到 squirrel"这个判断本身可能仍然成立，只是路径描述错了。

**处理方式**：这一条不用等 Phase 1 结束，**现在就可以顺手改**——先把第 88 行的路径改成实际路径（`internal/repository/monitoring/performance_log_repository.go`、`internal/repository/chat/chat_repository.go`）。如果 Phase 1 事务改造（A.1，需要给每个 Model 加 `WithSession` 方法）顺带把这两个文件也迁移到 squirrel 了，第 88 行这条"已知例外"注记就可以整段删除；如果没顺带迁移，保留这条但路径要先改对。

## 2. Phase-3-end 同步批次（架构性重写，范围大，现在只描述范围不写内容）

**这一批次的实质性内容依赖 Part B/C 的文档（`15-service-boundaries.md` 到 `21-cd-and-deployment.md`）先写完、Phase 2/3 实际落地后才能确定最终形态，本篇现在只描述"届时要改哪些段落"，不预写替换文本——预写的话大概率会和实际落地细节对不上，不如到时候照着真实代码状态写。**

届时（Phase 3 结束）需要重写的范围：

1. **`AGENTS.md` 第 28 行"关键目录"整段、`.cursor/rules/10-go-code-style.mdc` 第 11-29 行"目录结构"整段**：`internal/{handler,logic}/<domain>/<module>/` 这套描述不再匹配代码库真实结构——届时主结构是 `services/<name>/internal/...`（每个 RPC 服务各自一套 handler/logic/repository/model/domain）+ 一个显著变薄的、只服务 gateway 自身关切的 `internal/`（详见 B.3 的目标目录结构）。这不是小修，是整段替换。
2. **`AGENTS.md` 第 71 行"目录分层"、`.cursor/rules/10-go-code-style.mdc` 相应段落**：需要补充事务/领域服务分层标准（对应 `01-architecture-target.md`/`02-transactions-and-uow.md` 落地后的最终版本）、Wire 中间件方案（`03-wire-and-middleware.md`）、RPC 调用约定（`16-rpc-conventions.md`）、异步事件约定（`17-async-eventing.md`）——这些现在都还没定稿到"可以写进规则文件"的程度。
3. **`AGENTS.md` 第 2.1 节"新增模块脚手架"**：`generate-sql.sh` 的域→服务映射表（B.2 提到的固定映射：iam/content/chat/task/sdk）落地后，脚手架流程会多一步"确定新表归属哪个服务的 `db/services/<service>/<module>/`"，需要补充说明。
4. **`AGENTS.md` 第 6 节"何时必须停下来问用户"**：是否要把 `10-dev-execution-and-review-points.md` 的开发期执行策略正式合并进来，是 Phase 3 收尾时要向用户提出的决定点（见 `10-dev-execution-and-review-points.md` 末尾），不是自动发生的——即使决定合并，具体怎么措辞也要到那时候再定。
5. **`.cursor/rules/00-workflow.mdc` 第 20-55 行"标准开发流程"**：如果 Phase 2 新增了 `generate-rpc.sh`、Phase 3 新增了 `generate-swagger.sh`，且这两个脚本会被纳入"新增模块"的标准流程（比如新增一个需要独立 RPC 服务的模块时要走 `generate-rpc.sh`），这一节的步骤编号需要相应插入新步骤。

## 3. 执行方式

- **1.1/1.2 两条**：等 Phase 1 Week 3（"简单域 + 全仓库扫尾"）实际完成、`svcCtx.Domain` 真正成为 100% 的访问方式之后再改——改早了会出现"规则文件说的和当前代码实际状态不一致"的新问题（比如 Week 1 就把"旧代码可内联"删了，但 Week 1 时大部分 logic 文件其实还没迁移完，删了这句反而让规则和现实脱节）。
- **1.3 条**：现在就可以改（纯粹是修正一个已经错误的路径引用，不依赖本轮任何未完成的工作）。
- **第 2 节的五条**：Phase 3 对应文档（`15`-`21` 号）定稿、且 Phase 2/3 代码真正落地后才动手，本篇的第 2 节内容届时作为"要改哪几处"的检查清单使用，不是最终替换文本。

## 完成的定义

- Phase 1 结束时：`AGENTS.md` 第 28/71/72 行与 `.cursor/rules/10-go-code-style.mdc` 对应段落已按 1.1/1.2 更新，`.cursor/rules/10-go-code-style.mdc` 第 88 行路径已按 1.3 更正或整段删除，两处文件改动内容一致（`CLAUDE.md` 要求的同步约束）。
- Phase 3 结束时：第 2 节列出的五处全部核对并重写完成，且改完之后通读一遍 `AGENTS.md`/`.cursor/rules/*.mdc`，确认没有遗留任何指向 Phase 1 之前目录结构的死引用。
