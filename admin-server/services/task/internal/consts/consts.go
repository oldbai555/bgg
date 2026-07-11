// Package consts 复制自 internal/consts/consts.go 里 task 相关的常量。
// 按 16-rpc-conventions.md 第 6 节的既定策略：直接复制到各服务自己的 internal/consts，
// 不做成共享包（量很小，维护成本可忽略），后续两边如需变更需要各自同步改。
package consts

const (
	// TaskTypeExcelExport 异步导出Excel
	TaskTypeExcelExport int64 = 1
)

// SuperAdminUserID 是 Phase 1 种子数据里超级管理员账号的固定 ID（db/services/iam/user/
// init_user.sql），TaskCancel 的"管理员可以取消任何人的任务"判断沿用了这个既有约定（原
// internal/logic/task/task/task_cancel_logic.go 搬迁前就是这么写的，不是本轮新引入）。
// 这不是真正的 RBAC 校验，只是历史遗留的硬编码判断——真正的修法是 gateway 侧按权限
// （如 task:cancel:any）判断后把结果传给 RPC，或 task-rpc 反过来回调 iam-rpc 查角色，
// 两种都需要额外设计，留给后续会话评估是否要做，不在本轮 task-rpc 拆分范围内展开。
const SuperAdminUserID uint64 = 1

const (
	// TaskStatusPending 未开始
	TaskStatusPending int64 = 1
	// TaskStatusRunning 进行中
	TaskStatusRunning int64 = 2
	// TaskStatusCompleted 已完成
	TaskStatusCompleted int64 = 3
	// TaskStatusFailed 失败
	TaskStatusFailed int64 = 4
)

const (
	// TaskExecutionTypeSync 同步执行
	TaskExecutionTypeSync int64 = 1
	// TaskExecutionTypeAsync 异步执行
	TaskExecutionTypeAsync int64 = 2
)

const (
	// NotificationSourceTypeTask 任务通知来源类型
	NotificationSourceTypeTask = "task"
)

const (
	// RedisTaskLockPrefix 任务锁前缀（用于分布式锁）
	RedisTaskLockPrefix = "task:lock:"
)

const (
	// TaskNotificationTitleRunning 任务执行中通知标题
	TaskNotificationTitleRunning = "任务执行中"
	// TaskNotificationTitleCompleted 任务完成通知标题
	TaskNotificationTitleCompleted = "任务执行完成"
	// TaskNotificationTitleFailed 任务失败通知标题
	TaskNotificationTitleFailed = "任务执行失败"
)

const (
	// TaskDefaultScanInterval 默认扫描间隔（秒）
	TaskDefaultScanInterval = 5
	// TaskDefaultMaxConcurrent 默认最大并发执行数量
	TaskDefaultMaxConcurrent = 10
	// TaskDefaultBatchSize 默认批次大小
	TaskDefaultBatchSize = 100
	// TaskDefaultTaskTimeout 默认任务超时时间（秒，30分钟）
	TaskDefaultTaskTimeout = 1800
	// TaskDefaultLockTimeout 默认锁超时时间（秒）
	TaskDefaultLockTimeout = 30
)

const (
	// WSTaskProgress 任务进度消息类型
	WSTaskProgress = "task_progress"
	// WSNotification 通知消息类型
	WSNotification = "notification"
)

const (
	// TaskNotificationLevelInfo 信息级别
	TaskNotificationLevelInfo = "info"
	// TaskNotificationLevelSuccess 成功级别
	TaskNotificationLevelSuccess = "success"
	// TaskNotificationLevelError 错误级别
	TaskNotificationLevelError = "error"
)

// 任务导出模块常量（对应各业务模块的标识），透传给 TaskCallback.FetchExportData 的 module 字段。
const (
	TaskModuleOperationLog   = "operation_log"
	TaskModuleAuditLog       = "audit_log"
	TaskModuleLoginLog       = "login_log"
	TaskModuleSdkCallLog     = "sdk_call_log"
	TaskModulePerformanceLog = "performance_log"
)

// 导出任务筛选条件的 key（ExcelExportParams.Filters 的 key），透传进 filters_json。
const (
	TaskFilterUserId        = "userId"
	TaskFilterUsername      = "username"
	TaskFilterOperationType = "operationType"
	TaskFilterOperationObj  = "operationObject"
	TaskFilterMethod        = "method"
	TaskFilterPath          = "path"
	TaskFilterStartTime     = "startTime"
	TaskFilterEndTime       = "endTime"
	TaskFilterAuditType     = "auditType"
	TaskFilterAuditObject   = "auditObject"
	TaskFilterStatus        = "status"
	TaskFilterIsSlow        = "isSlow"
	TaskFilterStatusCode    = "statusCode"
	TaskFilterSdkKeyId      = "sdkKeyId"
	TaskFilterApiCode       = "apiCode"
	TaskFilterRespCode      = "respCode"
	TaskFilterIP            = "ip"
)

// UploadDir 上传文件存储目录（task-rpc 本地生成 CSV 的落盘目录，和 gateway 的
// PathFileUploads 通过 docker-compose 共享卷保持同一份物理文件，见 17-async-eventing.md）。
const UploadDir = "./uploads"

// PathFileUploads 文件下载代理路径前缀（gateway 侧的路由，用于拼 access_url）。
const PathFileUploads = "/api/v1/files/uploads"
