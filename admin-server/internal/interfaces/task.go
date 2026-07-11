package interfaces

import (
	"context"

	"postapocgame/admin-server/internal/model/task"
)

// TaskExecutor 任务执行器接口（面向具体业务的执行实现）
// 说明：不同业务类型的任务需要实现此接口，例如：
// - Excel导出任务执行器
// - 邮件通知任务执行器
// - 数据备份任务执行器
//
//go:generate mockery --name=TaskExecutor --output=../mocks/interfaces --outpkg=interfaces_mocks
type TaskExecutor interface {
	// GetType 获取任务类型（对应字典 task_type 的 value）
	GetType() int

	// Execute 执行任务
	// ctx: 上下文
	// task: 任务信息
	// paramsJSON: 任务参数（JSON字符串，对应TaskParamsReq的序列化结果）
	// 返回：任务结果（JSON字符串，对应TaskResultResp的序列化结果）和错误信息
	Execute(ctx context.Context, task *task.AdminTask, paramsJSON string) (string, error)
}

// AsyncTaskBackend 异步任务后端接口（面向调度引擎 / 中间件）
// 说明：通过该接口可以对接不同的任务中间件，例如：
// - 内置 go-zero 调度器（当前方案）
// - xxl-job / quartz / k8s CronJob 等第三方中间件
type AsyncTaskBackend interface {
	// Submit 提交一个任务到后端（返回后端任务ID，可与 admin_task.id 进行映射）
	Submit(ctx context.Context, task *task.AdminTask) (backendTaskID string, err error)

	// Cancel 取消一个任务（如果后端支持）
	Cancel(ctx context.Context, task *task.AdminTask) error

	// SyncStatus 同步任务状态（从后端获取任务状态，更新到 admin_task 表）
	// 说明：对于外部中间件（如 xxl-job），需要定期调用此方法同步状态
	SyncStatus(ctx context.Context, task *task.AdminTask) error

	// Start 启动后端（例如：启动调度器、连接中间件等）
	Start(ctx context.Context) error

	// Stop 停止后端（例如：停止调度器、断开连接等）
	Stop(ctx context.Context) error
}
