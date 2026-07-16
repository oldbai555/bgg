package interfaces

import (
	"context"

	"postapocgame/admin-server/services/task/internal/model/task"
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
