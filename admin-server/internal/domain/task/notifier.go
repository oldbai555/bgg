package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/hub"

	"github.com/zeromicro/go-zero/core/logx"
	systemmodel "postapocgame/admin-server/internal/model/system"
	taskmodel "postapocgame/admin-server/internal/model/task"
	"postapocgame/admin-server/internal/repository"
	systemrepo "postapocgame/admin-server/internal/repository/system"
)

// TaskNotifier 任务通知器
type TaskNotifier struct {
	repo    *repository.Repository
	chatHub *hub.ChatHub
}

// NewTaskNotifier 创建任务通知器
func NewTaskNotifier(repo *repository.Repository, chatHub *hub.ChatHub) *TaskNotifier {
	return &TaskNotifier{
		repo:    repo,
		chatHub: chatHub,
	}
}

// NotifyTaskStatusChange 通知任务状态变更
func (n *TaskNotifier) NotifyTaskStatusChange(ctx context.Context, task *taskmodel.AdminTask) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("任务通知发生 panic: %v, taskId=%d", r, task.Id)
		}
	}()

	// 1. 创建通知记录（持久化到数据库）
	n.createNotificationRecord(ctx, task)

	// 2. 通过WebSocket发送实时消息
	n.sendWebSocketMessage(ctx, task)
}

// createNotificationRecord 创建通知记录
func (n *TaskNotifier) createNotificationRecord(ctx context.Context, task *taskmodel.AdminTask) {
	notificationRepo := systemrepo.NewNotificationRepository(n.repo)

	var title, content string
	switch task.Status {
	case consts.TaskStatusRunning: // 进行中
		title = consts.TaskNotificationTitleRunning
		content = fmt.Sprintf("任务「%s」正在执行中...", task.Name)
	case consts.TaskStatusCompleted: // 已完成
		title = consts.TaskNotificationTitleCompleted
		content = fmt.Sprintf("任务「%s」执行完成", task.Name)
	case consts.TaskStatusFailed: // 失败
		title = consts.TaskNotificationTitleFailed
		content = fmt.Sprintf("任务「%s」执行失败：%s", task.Name, task.ErrorMessage)
	default:
		return // 其他状态不创建通知
	}

	now := time.Now().Unix()
	notification := &systemmodel.AdminNotification{
		UserId:     task.UserId,
		SourceType: consts.NotificationSourceTypeTask, // 使用常量
		SourceId:   task.Id,
		Title:      title,
		Content:    content,
		ReadStatus: 1, // 未读（字典值：1=未读，2=已读）
		ReadAt:     0,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  0,
	}

	if err := notificationRepo.Create(ctx, notification); err != nil {
		logx.Errorf("创建任务通知记录失败: taskId=%d, userId=%d, error: %v", task.Id, task.UserId, err)
	}
}

// sendWebSocketMessage 发送WebSocket实时消息
func (n *TaskNotifier) sendWebSocketMessage(ctx context.Context, task *taskmodel.AdminTask) {
	if n.chatHub == nil {
		return // ChatHub未初始化，跳过WebSocket通知
	}

	// 构建WebSocket消息
	var msgType string
	var level string
	var message string

	switch task.Status {
	case consts.TaskStatusRunning: // 进行中
		msgType = consts.WSTaskProgress
		level = consts.TaskNotificationLevelInfo
		message = fmt.Sprintf("任务「%s」正在执行中...", task.Name)
	case consts.TaskStatusCompleted: // 已完成
		msgType = consts.WSNotification
		level = consts.TaskNotificationLevelSuccess
		message = fmt.Sprintf("任务「%s」执行完成", task.Name)
	case consts.TaskStatusFailed: // 失败
		msgType = consts.WSNotification
		level = consts.TaskNotificationLevelError
		message = fmt.Sprintf("任务「%s」执行失败：%s", task.Name, task.ErrorMessage)
	default:
		return // 其他状态不发送WebSocket消息
	}

	// 使用ChatHub的ChatMessage结构体
	chatMsg := &hub.ChatMessage{
		Type:      msgType,
		TaskID:    fmt.Sprintf("%d", task.Id), // ChatMessage中TaskID是string类型
		TaskName:  task.Name,
		Status:    fmt.Sprintf("%d", task.Status), // 任务状态（字符串）
		Title:     fmt.Sprintf("任务「%s」", task.Name),
		Level:     level,
		Content:   message,
		CreatedAt: time.Now().Unix(),
	}

	// 序列化为JSON
	messageBytes, err := json.Marshal(chatMsg)
	if err != nil {
		logx.Errorf("任务通知消息序列化失败: taskId=%d, error: %v", task.Id, err)
		return
	}

	// 通过ChatHub发送给任务创建用户
	sent := n.chatHub.SendToUser(task.UserId, messageBytes)
	if !sent {
		logx.Infof("任务通知WebSocket发送失败（用户可能不在线）: taskId=%d, userId=%d", task.Id, task.UserId)
	}
}
