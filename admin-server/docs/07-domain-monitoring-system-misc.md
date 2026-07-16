# 07 — Monitoring / System / Misc 三域改造 + 全仓库扫尾（Phase 1 Week 3 收尾）

> 本文档是可直接执行的任务说明。这是 Phase 1 Week 3 的最后一份域改造文档，完成后紧接着做"Day 5 全仓库扫尾"（第 4 节），把 04/05/06/07 四份文档遗留的、真正无需领域服务的普通 logic 文件统一改成直接走 `svcCtx.Domain.X.Y`，收尾整个 Phase 1 的 DI 迁移。

## 0. 前置依赖

- [`01-architecture-target.md`](./01-architecture-target.md) —— 特别是"monitoring/system/misc 不单独建 `internal/domain/<x>/`，代码并入 `internal/domain/iam`（或紧邻的 `platform` 子包）"这条分组约定。
- [`02-transactions-and-uow.md`](./02-transactions-and-uow.md)
- [`03-wire-and-middleware.md`](./03-wire-and-middleware.md)
- [`04-domain-iam-chat.md`](./04-domain-iam-chat.md)、[`05-domain-task.md`](./05-domain-task.md)、[`06-domain-blog-video-sdk.md`](./06-domain-blog-video-sdk.md) —— 本文档的"Day 5 扫尾"任务依赖这三份文档已经落地的领域服务作为参照基准。

## 1. Monitoring 域：确认"薄"，无代码改动

`internal/logic/monitoring/**` 共 **17 个文件**（`audit_log` 3、`login_log` 4、`metric` 2、`metric_admin` 1、`monitor` 2、`operation_log` 3、`performance_log` 2）。逐个 grep 多仓储调用和写操作，结果：

- 全部 17 个文件里没有一个出现 ≥2 个不同仓储构造调用（"monitoring multi-repo" 检查结果为空）。
- 唯一的写操作集中在两类：
  1. `metric_report_logic.go` 的 `MetricRepository.UpsertDailyStats(...)`（`internal/logic/monitoring/metric/metric_report_logic.go:127`）——单表 `upsert`，一条 SQL，不需要事务。
  2. `audit_log_export_logic.go`/`login_log_export_logic.go`/`operation_log_export_logic.go`/`performance_log_export_logic.go` 四个导出文件——它们不直接写日志表，而是往 `admin_task` 表插一行导出任务（跨到 Task 域的 `taskrepo`），单表写，语义是"发起一个任务"，不是需要原子性保护的组合写。这四个文件正是计划文档 Part B.2 里点名的"4 个 `monitoring/*_export_logic.go`"跨域 import，**只读/单写、不满足领域服务门槛**，处理方式同 06 文档 2.4/4.1 节——加 TODO 注释标记 Phase 2 待办，本文档不修。
- 日志表本身的写入（`admin_operation_log`/`admin_performance_log`）发生在中间件（`OperationLogMiddleware`/`PerformanceMiddleware`），不在 `internal/logic` 里，不属于本文档范围（中间件改造属于 03 文档）。

**结论：Monitoring 域按设计就是薄的（"读多写少、写的地方语义都是单表"），确认符合计划文档的判断，不需要任何领域服务，不做代码改动。**

## 2. System 域：审计结论修正计划文档的猜测

`internal/logic/system/**` 共 **30 个文件**（`config` 5、`dict` 2、`dict_item` 4、`dict_type` 4、`file` 6、`notice` 4、`notification` 5），文件数和计划文档估算的"30"完全吻合。多仓储调用统计：

| 文件 | 仓储数 | 实际操作 | 判断 |
|---|---|---|---|
| `file/file_upload_logic.go` | 3 | 读字典(`DictType`+`DictItem`) 取 baseURL，只写 `admin_file` 一张表 | 单表写，不需要领域服务 |
| `dict_item/dict_item_create\|delete\|update_logic.go` | 2 | 先读 `DictType` 校验存在性，只写 `admin_dict_item` 一张表 | 单表写，不需要领域服务 |
| `dict_type/dict_type_delete_logic.go` | 2 | 先读 `DictItem` 校验没有子项，只写 `admin_dict_type` 一张表 | 单表写，不需要领域服务 |
| `dict/dict_get_logic.go`、`dict/dict_batch_get_logic.go` | 2 | 纯读聚合 | 保持薄 logic |
| `notice/notice_create_logic.go`、`notice/notice_update_logic.go` | 3（含跨域 `iamrepo.NewUserRepository`） | 写 `admin_notice`（同步）+ 为全体用户批量写 `admin_notification`（**已经是** `go l.createNotificationsForAllUsers(...)` 异步 goroutine，内部已经用 `userRepo.FindChunk` 分批处理，失败只 `logx.Errorf`，见 `internal/logic/system/notice/notice_create_logic.go:86,103-167`） | **已经是正确的异步尽力而为写法，不需要 Transact** |

**这里需要明确修正计划文档的猜测**：计划文档原文写"system 的 notice-create 触及 dict+notice 表是唯一真实的多表案例"——实际读代码后发现这个描述有两处不准：① `notice_create_logic.go` 根本不碰字典表，跨的是 `admin_notice` + `admin_notification`；② 这次跨表写**已经是异步、已经用了 `FindChunk` 分批**（很可能是上一轮或更早的修复已经顺手做了，和 04 文档里 `user_create_logic.go` 至今仍在用 `FindPage(1, 10000, "")` 全表拉取形成对比）。也就是说 **System 域实际上没有遗留的多表写事务缺口**，唯一称得上"多表"的场景已经用异步最终一致的方式处理妥当，不需要再包 `Transact`。

`notice_create_logic.go`/`notice_update_logic.go` 里 `iamrepo.NewUserRepository` 这处跨域 import 和 06 文档 2.4 节的处理方式一致：只读遍历，加 TODO 注释标记 Phase 2 待办，本文档不修。

**结论：System 域没有需要领域服务的多表写方法，不做代码改动，只需要在 Day 5 扫尾时把 30 个文件里没有跨域 import 问题的部分（即除 `notice/*` 之外的 26 个）纳入统一的直调改造。**

## 3. Misc 域：确认最薄

`internal/logic/misc/**` 共 **10 个文件**（`daily_short_sentence` 4、`demo` 4、`ping` 1、`public` 1），和计划文档估算完全一致。全部是单表 CRUD（`demo`/`daily_short_sentence`）或无状态端点（`ping_logic.go` 健康检查、`public_dict_get_logic.go` 只读字典）。grep 多仓储调用结果为空。**确认是全仓库最薄的域，不需要任何领域服务改造。**

## 4. Day 5 扫尾：`l.svcCtx.Repository` → `svcCtx.Domain.X.Y` 全仓库收口

### 4.1 现状基线（本文档写作时实测，不是估算）

```
$ grep -rl "l.svcCtx.Repository" internal/logic | wc -l
161
$ grep -rl "svcCtx.Domain" internal/logic | wc -l
1
```

161 个 logic 文件里目前只有 1 个真正用上了 Wire 装配好的 `registry.Domain`，其余全部靠方法体内 `xxxrepo.NewXxxRepository(l.svcCtx.Repository)` 现场构造——这正是计划文档开篇提到的"`registry.Domain` 聚合已经搭好但 161 个 logic 文件里只有 1 个用上"的具体数字来源，现已核实无误。

### 4.2 04/05/06 落地后预期的变化

- 04 文档：`user_create_logic.go`、`role_permission_update_logic.go`、`permission_menu_update_logic.go`、`user_role_update_logic.go`、`permission_api_update_logic.go` 共 5 个文件改为调领域服务（`svcCtx.Domain.IAM.UserService`/`svcCtx.Domain.IAM.RBAC`），不再是"简单单仓储直调"，**不纳入本节的机械替换范围**（它们已经在 04 文档里改好了）。
- 06 文档：`blog_article_create_logic.go`、`blog_article_update_logic.go`、`blog_article_audit_logic.go`、`blog_article_audit_unpublish_logic.go`、`sdk_api_key_bind_save_logic.go` 共 5 个文件同理，改为调 `svcCtx.Domain.Blog.ArticleService`/`svcCtx.Domain.SDK.Service`（**不是** `Domain.Content.BlogArticle`——`registry.Domain` 顶层没有 `Content` 字段，`internal/domain/content` 只是 Go 包目录名，见 06 文档第 5 节），**不纳入本节**。
- 05 文档：Task 域的 3 个文件（`scheduler.go`/`notifier.go`/`excel_export_executor.go`）本来就在 `internal/domain/task/` 里，不是 `internal/logic` 下的文件，不计入 161 这个基数。

也就是说 161 里大约有 10 个文件会在 04/06 落地后从"直调仓储"变成"调领域服务"，**剩下约 151 个文件是本节要处理的"合法的简单单仓储/多仓储只读直调"**。

### 4.3 机械替换规则

对剩下的每个文件：

1. 找出方法体内所有 `<pkg>repo.New<Xxx>Repository(l.svcCtx.Repository)` 调用。
2. 对照 `internal/repository/registry/domain.go` 里 `NewDomain` 函数的字段映射（例如 `iamrepo.NewUserRepository(repo)` 对应 `IAMDomain.User`、`systemrepo.NewNoticeRepository(repo)` 对应 `SystemDomain.Notice`），把调用点替换成 `l.svcCtx.Domain.<Domain>.<Field>`，删除局部变量声明和 import。
3. **例外，不要替换**：
   - `l.svcCtx.Repository.BusinessCache` —— 这是缓存工具，不是仓储，`registry.Domain` 没有对应字段，本次实测有 **18 个文件**用到（`grep -rl "svcCtx.Repository.BusinessCache" internal/logic | wc -l`），保持原样。
   - 任何直接访问 `l.svcCtx.Repository.<XxxModel>`（goctl 生成的 Model，不经过 Repository 接口封装）的写法——本次审计发现 `blog_article_audit_logic.go:82` 有一处 `l.svcCtx.Repository.BlogArticleModel.Update(...)`，这类直连 Model 的写法本身就是技术债，但**修它属于"引入领域服务"的范畴（06 文档已经处理了这一个）**，不属于本节"纯换路径不改行为"的机械替换，如果扫尾时又发现类似的直连 Model 调用，记录下来但不要顺手改，避免扫尾任务范围失控。
4. 替换后跑 `go build ./...`，每处理完一个域（建议按 monitoring → system → misc → 剩余的 iam/blog/video/sdk/chat 只读文件顺序）跑一次编译，不要攒到最后一次性改完再编译。

### 4.4 验证

```bash
# 替换完成后，重新统计：
grep -rl "l.svcCtx.Repository" internal/logic | wc -l
```

预期剩余匹配只包括：① 18 个 `BusinessCache` 相关文件（第 4.3 节例外 1）；② 尚未处理的、本节明确排除在外的极少数直连 Model 文件（如果扫尾时发现新的，记录到 `11-descoped.md`，不在本轮处理）。如果剩余数字明显大于"18 + 少量已知例外"，说明机械替换没做完，回到 4.3 节继续。

```bash
grep -rl "svcCtx.Domain" internal/logic | wc -l
```

预期从 1 涨到 150+（约等于 161 - 10 个领域服务文件 + 本节新增的直调文件数，具体数字以实际替换结果为准，不强求刚好对上估算值）。

## 5. 非目标

- 不给 Monitoring/System/Misc 任何一个域引入 `internal/domain/<x>/` 包——按 01 文档的分组约定，如果这三个域将来真的需要领域服务，代码应该写进 `internal/domain/iam/`（或紧邻的 `platform` 子包），不新建目录。本轮审计结论是三个域都不需要，这条约定暂时用不上，但写文档时仍要明确，避免执行者按"每个域都要有 domain 包"的惯性思维新建空目录。
- 不修 `notice/*`、`monitoring/*_export_logic.go` 的跨域 import（原因见第 1、2 节，统一记录为 Phase 2 待办）。
- Day 5 扫尾不改任何业务行为，纯路径替换；如果替换过程中发现某个"看似简单"的文件其实有隐藏的多表写（理论上第 1-3 节已经全量审计过，不应该再有漏网之鱼，但如果真的发现了），停下来，不要在扫尾任务里顺手加事务保护，记录下来单独处理。

## 6. 完成的定义

1. `go build ./...` 通过。
2. `grep -rl "l.svcCtx.Repository" internal/logic | wc -l` 的结果和第 4.4 节预期一致（约 18 + 少量已知例外），并且这份最终数字和例外清单写进 `progress.md`。
3. Monitoring/System/Misc 三个域本身没有新增文件（因为审计结论是都不需要领域服务），`internal/domain/` 目录下没有出现 `monitoring/`、`system/`、`misc/` 子目录。
4. 人工冒烟测试：
   - 后台查看操作日志列表，触发一次操作日志导出，确认任务创建成功、能在任务中心看到进度。
   - 新建/编辑一条字典项，确认字典缓存清除逻辑仍然生效（改字典后前端下拉选项能刷新）。
   - 发布一条公告，确认全体用户异步收到未读通知（这个行为本次未改动，验证的是"扫尾没有破坏原有功能"）。
   - 访问 `/api/v1/misc/ping` 确认健康检查接口未受影响。
