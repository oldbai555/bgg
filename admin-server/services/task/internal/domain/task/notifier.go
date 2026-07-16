// Package task 从 internal/domain/task/notifier.go 改造而来。原实现直接持有
// systemrepo.NotificationRepository（写 admin_notification，物理属于 iam）和 *hub.ChatHub
// （WebSocket 推送，物理属于 chat），这两处都是 task-rpc 拆分后拿不到的跨域依赖，且原实现
// 两处失败都只 logx.Errorf、不影响任务主流程——完全符合 17-async-eventing.md 的判断规则
// "现在代码里子操作失败只是记日志、不影响主流程的 → 异步 Streams"（对应计划里的"发现 4"）。
//
// 改造成发布 stream:task.notification 事件，消费者（当前阶段：单体内新增的一个消费者，
// 因为 iam-rpc/chat-rpc 都还没真正拆分成独立进程，是一个跨两个未来域的临时合并消费者）
// 负责"写 admin_notification 记录 + 推 WS"两件事。幂等/死信策略与 stream:chat.user.created
// 同构，见 17-async-eventing.md 第 2.3/2.4 节。
package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"postapocgame/admin-server/services/task/internal/consts"
	taskmodel "postapocgame/admin-server/services/task/internal/model/task"
)

// StreamTaskNotification 是 task-rpc 发布任务状态变更事件的 Redis Stream 名称。
const StreamTaskNotification = "stream:task.notification"

// TaskNotificationEvent 是发布到 StreamTaskNotification 的事件体（JSON 编码后存进
// payload 字段），消费者据此决定通知标题/内容/WS 消息级别，字段命名与内容直接对应
// 原 notifier.go 里 createNotificationRecord/sendWebSocketMessage 两段逻辑用到的数据。
type TaskNotificationEvent struct {
	TaskID       uint64 `json:"taskId"`
	TaskName     string `json:"taskName"`
	UserID       uint64 `json:"userId"`
	Status       int64  `json:"status"` // consts.TaskStatus*
	ErrorMessage string `json:"errorMessage,omitempty"`
	SourceType   string `json:"sourceType"` // 固定 consts.NotificationSourceTypeTask
}

// TaskNotifier 任务通知器
type TaskNotifier struct {
	redis *redis.Redis
}

// NewTaskNotifier 创建任务通知器
func NewTaskNotifier(rds *redis.Redis) *TaskNotifier {
	return &TaskNotifier{redis: rds}
}

// NotifyTaskStatusChange 通知任务状态变更：发布一条 stream:task.notification 事件。
// 只在 Running/Completed/Failed 三种状态发布（与原实现的过滤规则一致，其他状态不通知）。
func (n *TaskNotifier) NotifyTaskStatusChange(ctx context.Context, task *taskmodel.AdminTask) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("任务通知发生 panic: %v, taskId=%d", r, task.Id)
		}
	}()

	switch task.Status {
	case consts.TaskStatusRunning, consts.TaskStatusCompleted, consts.TaskStatusFailed:
	default:
		return // 其他状态不通知
	}

	event := TaskNotificationEvent{
		TaskID:       task.Id,
		TaskName:     task.Name,
		UserID:       task.UserId,
		Status:       task.Status,
		ErrorMessage: task.ErrorMessage,
		SourceType:   consts.NotificationSourceTypeTask,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		logx.Errorf("任务通知事件序列化失败: taskId=%d, error: %v", task.Id, err)
		return
	}

	// 发布失败同样只记日志，不影响任务主流程——与被替换前的直接写库/推送失败处理语义一致。
	if _, err := n.redis.XAddCtx(ctx, StreamTaskNotification, false, "*", []string{"payload", string(payload)}); err != nil {
		logx.Errorf("发布 %s 事件失败: taskId=%d, error: %v", StreamTaskNotification, task.Id, err)
		return
	}

	logx.Infof("已发布任务通知事件: taskId=%d, status=%d, %s", task.Id, task.Status, fmt.Sprintf("stream=%s", StreamTaskNotification))
}
