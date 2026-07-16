# 05 — Task 域：对齐 Transact + 补测试（Phase 1 Week 2）

> 本文档是可直接执行的任务说明。执行者在改动前应完整阅读一遍，改完跑 `go build ./...` 和相关测试。

## 0. 前置依赖

- [`01-architecture-target.md`](./01-architecture-target.md)
- [`02-transactions-and-uow.md`](./02-transactions-and-uow.md) —— `repository.Repository.Transact(ctx, fn)`
- [`03-wire-and-middleware.md`](./03-wire-and-middleware.md)

## 1. 为什么这份文档比 04/06/07 小

`internal/domain/task/` 是全仓库唯二已经有真正领域服务层的域之一（另一个是 IAM 的 `PermissionResolver`）。上一轮 DDD-lite 重构已经把 `internal/task/*` 整体搬到了 `internal/domain/task/`，`internal/interfaces.TaskExecutor`/`AsyncTaskBackend` 两个接口从一开始就是按"可插拔后端"设计的（这也是 Part B 里 `task-rpc` 第一个被拆出去的原因）。本文档要做的不是新建领域服务，而是**审计现有 3 个核心文件的写操作是否需要 `Transact` 保护，逐个确认或修正，然后补测试**。审计结论（见第 2 节）是：**当前没有发现真正需要事务保护的多表写场景**，这和 04/06/07 里发现的多个真实事务缺口不同——这正是本文档篇幅明显小于其他三份的原因。

## 2. 逐文件审计结果

### 2.1 `internal/domain/task/scheduler.go`

- `scanAsyncTasks`/`scanScheduledTasks` 直接用 `s.repo.DB` + squirrel 查询 `admin_task` 表（第 132-158、161-190 行），**读操作绕开了传进来的 `taskRepo` 参数**——`scanAndExecute` 里构造了 `taskRepo := taskrepo.NewTaskRepository(s.repo)`（第 98 行）并作为参数传给这两个方法，但方法体内根本没用这个参数，是上一轮重构遗留的"接口留着、实现没跟上"。这属于"squirrel 直连改造为通过 repository 调用"的范畴，**上一轮 DDD-lite 任务书（`docs/admin-server-ddd-refactor-prompt.md` 第 5 节）已经明确写了"不要顺带把 scheduler 内部直连 squirrel 的 SQL 改成走 repository 接口——这是可选的后续优化"，本轮同样不强制修**。如果顺手做，改动范围是：`scanAsyncTasks`/`scanScheduledTasks` 删掉未使用的 `taskRepo` 参数，或者真的通过 `taskRepo` 暴露的方法查询（需要先给 `TaskRepository` 接口加对应的条件查询方法）。**这是可选项，不是本文档的必做任务**，做与不做都不影响"完成的定义"。
- `executeTask` 里的写操作：`taskRepo.UpdateStatus(...)`（置为运行中）→ 执行 → `taskRepo.UpdateResult(...)`（置为完成/失败）。这两次写各自只操作 `admin_task` 单表的单行，且是任务生命周期里两个先后发生、语义上独立的状态迁移（不是同一个逻辑操作被拆成两次写），**不需要 Transact**——把它们包在一个事务里没有意义，中间还隔着可能耗时很久的 `executor.Execute(...)` 调用，事务不应该跨越这么长的执行窗口。
- `handleTaskError` 同理，单表单行写，不需要 Transact。
- `acquireLock`/`releaseLock` 是 Redis 操作，不涉及 MySQL 事务，不在本文档范围内（Redis 分布式锁的正确性已经由 `SETEX` 的原子性保证，计划文档 B.1 也确认"调度器现有的 Redis 锁已经是多副本安全的，不需要重新设计并发模型"）。

结论：`scheduler.go` **不需要任何 Transact 改造**。

### 2.2 `internal/domain/task/notifier.go`

- `NotifyTaskStatusChange` 依次调用 `createNotificationRecord`（写 `admin_notification` 表，跨到 System 域的 `systemrepo.NewNotificationRepository`）和 `sendWebSocketMessage`（走 `ChatHub`，内存态，不是 DB 写）。`createNotificationRecord` 内部失败只 `logx.Errorf`，不返回 error 给调用方（`createNotificationRecord`/`sendWebSocketMessage` 都是 `func(...)` 无返回值，第 49、87 行），任务状态本身已经在 `scheduler.go` 里提交过了——这正是计划文档 B.4 定义的"现在代码里子操作失败只是 `logx.Errorf` 记录、不影响主流程"的模式，**按规则应该保持异步尽力而为，不需要用 Transact 把它和任务状态更新绑在一起**（绑在一起反而是倒退：任务本身执行成功了，不该因为通知写入失败就回滚任务状态）。
- 这里有一处跨域 import（`systemrepo.NewNotificationRepository`，`notifier.go:16`）——Task 域直接持有 System 域仓储。这属于 Part B.2 里"没有 FK、只是 Go 代码里跨域直接 import"这一类问题，但因为它是**只读的通知写入、失败不影响主流程**，不满足"跨域调用应该走窄接口"的紧迫性（04 文档里 IAM→Chat 那条要treat 是因为 IAM 建用户是核心业务流程且 Chat 越界更严重）。本文档不要求为这一处引入 `systemdomain.Notifier` 之类的窄接口，标记为**已知技术债，留给 Phase 2 拆分 `task-rpc`/`iam-rpc` 时通过 `TaskCallback`/Streams 机制自然消解**（见计划文档 B.4），本轮不动。

结论：`notifier.go` **不需要 Transact 改造**，跨域 import 记录为已知项，不在本文档修复范围。

### 2.3 `internal/domain/task/executors/excel_export_executor.go`

- `Execute` 本身只读（查询 4 种日志表之一），不涉及写。
- `generateCSVFile` 里有真正的"两种资源要保持一致"的场景：先在本地文件系统写 CSV 文件，再 `fileRepo.Create(ctx, &fileModel)` 写 `admin_file` 表；如果 DB 写入失败，代码已经手动 `os.Remove(fileSystemPath)` 做补偿删除（第 490-494 行）。这是**单表 DB 写 + 文件系统写**的组合，不是多表 DB 写，SQL 事务本来就管不到文件系统这一侧，现有的"失败就删文件"补偿逻辑已经是这种场景下正确的处理方式，**不需要也不能用 `Transact` 覆盖**（`Transact` 只能保证 MySQL 内部的原子性）。
- `getStorageBaseURL` 只读（查字典类型 + 字典项两次读），不涉及写。

结论：`excel_export_executor.go` **不需要 Transact 改造**。

## 3. 唯一要做的代码改动（可选，非阻塞）

如果顺手处理，把 `scanAsyncTasks`/`scanScheduledTasks` 未使用的 `taskRepo taskrepo.TaskRepository` 参数删掉（避免误导读者以为这两个方法真的在用它）：

```go
func (s *TaskScheduler) scanAsyncTasks(ctx context.Context) ([]taskmodel.AdminTask, error) {
	// 函数体不变，只删参数
}
```

调用点 `scanAndExecute` 里同步去掉 `taskRepo` 变量和两处传参。这是纯清理，不是本文档"完成的定义"的必要条件。

## 4. 测试任务（本文档的主要工作量）

零测试起点，按 [`08-testing-strategy.md`](./08-testing-strategy.md)（sqlmock，happy-path 优先，只在有真实分支的地方补）给三个文件各补测试：

### 4.1 `scheduler_test.go`

- `TestScanAsyncTasks`/`TestScanScheduledTasks`：mock `admin_task` 表的 `SELECT`，断言生成的 SQL 里 WHERE 条件包含 `execution_type=2`、`status=1`、`scheduled_at=0`（或 `>0 AND <=now`），断言 `ORDER BY`/`LIMIT` 符合 `consts.TaskDefaultBatchSize`。
- `TestExecuteTask_Success`：mock `acquireLock` 对应的 Redis 调用（用 `miniredis` 或直接 mock `*repository.Repository.Redis`，视 02 文档定的 Redis mock 约定而定）返回加锁成功 → `FindOne` 返回 `status=pending` → `UpdateStatus` 成功 → executor mock 返回成功结果 → `UpdateResult` 成功 → 断言 `notifier.NotifyTaskStatusChange` 被调用两次（运行中 + 完成）。
- `TestExecuteTask_LockHeld`：mock 加锁失败（锁已存在），断言 `UpdateStatus`/executor 都不会被调用（任务被跳过）。
- `TestExecuteTask_ExecutorNotFound`：`executors` map 里没有对应 `task.Type`，断言 `handleTaskError` 被调用且最终状态是 `TaskStatusFailed`。
- `TestExecuteTask_ExecutorPanic`：executor mock 直接 panic，断言 `recover()` 生效、任务被标记失败、不会导致整个调度器 goroutine 崩溃。

### 4.2 `notifier_test.go`

- `TestNotifyTaskStatusChange_Running/Completed/Failed`：分别断言 `createNotificationRecord` 生成的 `title`/`content` 符合 `consts.TaskNotificationTitleXxx` 常量，`sendWebSocketMessage` 组装的 `ChatMessage.Type` 符合 `consts.WSTaskProgress`/`consts.WSNotification`。
- `TestNotifyTaskStatusChange_NotificationCreateFails`：`notificationRepo.Create` mock 返回 error，断言方法本身不 panic、不 return error（因为签名就是无返回值)，只走 `logx.Errorf` 分支（可以通过检查 mock 的 `ExpectationsWereMet()` 确认这条 INSERT 确实被尝试过一次）。
- `TestSendWebSocketMessage_ChatHubNil`：`chatHub` 为 nil 时直接 return，不 panic。

### 4.3 `excel_export_executor_test.go`

- `TestExecute_UnsupportedModule`：`params.Module` 传一个不在 `switch` 分支里的值，断言返回 `"不支持的导出模块"` 错误。
- `TestExportOperationLog_Success`：mock `operationLogRepo.FindPage` 返回 2 条记录，断言生成的 CSV 内容行数、表头符合预期（可以用临时目录 + 读文件内容断言，不需要 mock 文件系统）。
- `TestGenerateCSVFile_DBWriteFailsCleansUpFile`：mock `fileRepo.Create` 返回 error，断言方法返回 error 且临时创建的 CSV 文件已被删除（`os.Stat` 返回 `IsNotExist`）——这是验证第 2.3 节提到的"补偿删除"逻辑确实生效，是本文件测试里价值最高的一个用例。
- `TestGenerateCSVFile_FileAlreadyExists`：mock `fileRepo.FindByName` 命中已存在记录，断言直接复用已有记录、不重复写 `admin_file`。

跳过测试：`calculateFileMD5`（纯工具函数，逻辑简单）、`getStorageBaseURL` 的字典缺失分支之外的其余分支可以合并到一个用例里测（不需要每个字段都单独开一个 test case）。

## 5. 非目标

- 不重写 `scheduler.go` 内部直连 squirrel 的查询逻辑（第 3 节的参数清理除外，且是可选的）。
- 不给 `TaskExecutor`/`AsyncTaskBackend` 接口新增方法或改签名——这两个接口是 Phase 2 `task-rpc` 拆分的既定契约，本轮只补测试不动接口。
- 不引入 `systemdomain.Notifier` 窄接口去修 `notifier.go` 的跨域 import（第 2.2 节已说明原因）。
- 不测试 goroutine 并发时序本身（`maxConcurrent`/`semaphore` 的并发正确性）——计划文档 `11-descoped.md` 明确"不做 WebSocket 聊天的并发/压力测试"，任务调度器的并发压测同样不在本轮范围。

## 6. 完成的定义

1. `go build ./...` 通过。
2. `go test ./internal/domain/task/... -v` 全部通过，覆盖第 4 节列出的用例。
3. 如果做了第 3 节的可选清理，确认 `scanAsyncTasks`/`scanScheduledTasks` 调用点参数同步更新，`go vet ./...` 无未使用变量告警。
4. 人工冒烟测试：本地起服务，触发一次 Excel 导出任务（调 `POST /api/v1/monitoring/operation_log/export` 或类似接口），确认几秒内任务状态从"待执行"变成"进行中"再到"已完成"，`admin_notification` 表里能查到对应记录，WebSocket 连接（如果前端已连上）能收到任务完成通知。
