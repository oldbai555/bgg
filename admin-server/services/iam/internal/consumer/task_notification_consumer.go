// Package consumer 消费 pkg/taskcallback（准确说是 services/task/internal/domain/task
// TaskNotifier）发布的 stream:task.notification 事件，负责"写 admin_notification 记录 +
// 推 WS"两件事——这是 17-async-eventing.md 计划要求的、原 internal/domain/task/notifier.go
// 直接持有 systemrepo.NotificationRepository + *hub.ChatHub 的那部分逻辑，现在通过 Streams
// 解耦，不再是进程内直接调用。
//
// 消费者组名 iam-chat-task-notify 刻意体现"这是一个跨两个未来域（iam 的通知表 + chat 的
// WebSocket 推送）的临时合并消费者"——iam-rpc/chat-rpc 真正拆分成独立进程后，这里要拆成
// 两个消费者分别处理，见 docs/progress.md 本轮条目。
package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"postapocgame/admin-server/services/iam/internal/consts"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/repository"
	systemrepo "postapocgame/admin-server/services/iam/internal/repository/system"
	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/chatclient"
)

const (
	// StreamTaskNotification 必须和 services/task/internal/domain/task/notifier.go 里的
	// 同名常量保持一致，两边各自维护一份（16-rpc-conventions.md 第 6 节"直接复制不共享"）。
	StreamTaskNotification = "stream:task.notification"
	streamDeadLetterSuffix = ".deadletter"
	consumerGroup          = "iam-rpc-task-notify"
	maxDeliveryAttempts    = 3
)

// taskNotificationEvent 与 services/task 侧 TaskNotificationEvent 字段一一对应。
type taskNotificationEvent struct {
	TaskID       uint64 `json:"taskId"`
	TaskName     string `json:"taskName"`
	UserID       uint64 `json:"userId"`
	Status       int64  `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	SourceType   string `json:"sourceType"`
}

// TaskNotificationConsumer 消费 stream:task.notification。
type TaskNotificationConsumer struct {
	redis    *redis.Redis
	repo     *repository.Repository
	chatRPC  chatclient.Chat // chat 域已拆分成独立服务，推 WS 改成回调 chat-rpc.PushToUser
	stopChan chan struct{}
	wg       sync.WaitGroup

	// consumerName 必须在同一个消费者组内唯一，否则多副本部署时会互相抢占/无法正确分担
	// 负载（Redis Stream 消费者组按 consumer name 追踪各自的 pending 消息）。用 hostname+PID
	// 拼出来，同机多进程/多副本容器（hostname 各不相同）都能保证唯一。
	consumerName string

	failuresMu sync.Mutex
	failures   map[string]int // messageID -> 已失败次数，仅进程内存，重启后清零（见包注释的简化说明）
}

// NewTaskNotificationConsumer 创建消费者。
func NewTaskNotificationConsumer(rds *redis.Redis, repo *repository.Repository, chatRPC chatclient.Chat) *TaskNotificationConsumer {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "unknown-host"
	}
	return &TaskNotificationConsumer{
		redis:        rds,
		repo:         repo,
		chatRPC:      chatRPC,
		stopChan:     make(chan struct{}),
		failures:     make(map[string]int),
		consumerName: fmt.Sprintf("%s-%s-%d", consumerGroup, hostname, os.Getpid()),
	}
}

// Start 启动消费者组（幂等，BUSYGROUP 错误直接忽略）并起后台 goroutine 开始消费。
func (c *TaskNotificationConsumer) Start() {
	if _, err := c.redis.XGroupCreateMkStreamCtx(context.Background(), StreamTaskNotification, consumerGroup, "$"); err != nil {
		// BUSYGROUP：消费者组已存在，是预期情况（重启进程后组本身已经在 Redis 里持久化），不是错误。
		logx.Infof("XGROUP CREATE %s %s: %v（已存在则忽略）", StreamTaskNotification, consumerGroup, err)
	}

	c.wg.Add(1)
	go c.run()
	logx.Infof("task 通知消费者已启动: stream=%s, group=%s", StreamTaskNotification, consumerGroup)
}

// Stop 停止消费者。
func (c *TaskNotificationConsumer) Stop() {
	close(c.stopChan)
	c.wg.Wait()
	logx.Infof("task 通知消费者已停止")
}

func (c *TaskNotificationConsumer) run() {
	defer c.wg.Done()

	node, err := redis.CreateBlockingNode(c.redis)
	if err != nil {
		logx.Errorf("task 通知消费者创建 blocking node 失败: %v", err)
		return
	}
	defer node.Close()

	for {
		select {
		case <-c.stopChan:
			return
		default:
		}

		ctx := context.Background()
		streams, err := c.redis.XReadGroupCtx(ctx, node, consumerGroup, c.consumerName, 10, 2*time.Second, false,
			StreamTaskNotification, ">")
		if err != nil {
			// redis.Nil / 超时是正常情况（没有新消息），不打日志刷屏；其他错误记日志、稍等重试。
			if err != redis.Nil {
				logx.Errorf("XREADGROUP %s 失败: %v", StreamTaskNotification, err)
				time.Sleep(time.Second)
			}
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				c.handleMessage(ctx, msg.ID, msg.Values)
			}
		}
	}
}

func (c *TaskNotificationConsumer) handleMessage(ctx context.Context, msgID string, values map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("task 通知消费发生 panic: msgID=%s, error: %v", msgID, r)
		}
	}()

	payload, _ := values["payload"].(string)
	var event taskNotificationEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		logx.Errorf("解析 task 通知事件失败，直接丢弃: msgID=%s, error: %v", msgID, err)
		c.ack(ctx, msgID)
		return
	}

	if err := c.process(ctx, event); err != nil {
		c.handleFailure(ctx, msgID, event, err)
		return
	}

	c.failuresMu.Lock()
	delete(c.failures, msgID)
	c.failuresMu.Unlock()
	c.ack(ctx, msgID)
}

// process 对应原 notifier.go 的 createNotificationRecord + sendWebSocketMessage 两步。
func (c *TaskNotificationConsumer) process(ctx context.Context, event taskNotificationEvent) error {
	title, content, ok := notificationText(event)
	if !ok {
		return nil // 其他状态不通知，直接视为处理成功
	}

	notificationRepo := systemrepo.NewNotificationRepository(c.repo)

	// 幂等：同一条 taskId+title 的通知已存在则跳过（同一条 Stream 消息被重复投递时，
	// 不会重复建通知），策略与 17-async-eventing.md 第 2.3 节的"插入前查是否已存在"一致。
	// 按 created_at DESC 取最近 20 条足够覆盖重复投递场景（redelivery 发生在原消息处理后
	// 不久，同一用户在这期间产生的其他通知数量有限）；这里只是防止 Streams 至少一次投递
	// 语义下的重复插入，不是强一致性保证，真正兜底幂等的是最终人工可查的 deadletter 流。
	existing, _, err := notificationRepo.FindPage(ctx, 1, 20, event.UserID, consts.NotificationSourceTypeTask, -1)
	if err == nil {
		for _, n := range existing {
			if n.SourceId == event.TaskID && n.Title == title {
				return c.pushWS(ctx, event, title, content)
			}
		}
	}

	now := time.Now().Unix()
	notification := &systemmodel.AdminNotification{
		UserId:     event.UserID,
		SourceType: consts.NotificationSourceTypeTask,
		SourceId:   event.TaskID,
		Title:      title,
		Content:    content,
		ReadStatus: 1, // 未读（字典值：1=未读，2=已读）
		ReadAt:     0,
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  0,
	}
	if err := notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("创建任务通知记录失败: %w", err)
	}

	return c.pushWS(ctx, event, title, content)
}

// wsChatMessage 是 WS wire 格式（JSON）的本地副本，字段名/tag 与
// services/chat/internal/hub.ChatMessage 逐一对齐——chat 域拆分后那个类型定义在 chat-rpc
// 自己的 internal/ 下不能跨服务 import，按 16-rpc-conventions.md 第 6 节"直接复制不共享"
// 的既定策略，gateway 侧这里维护一份同形状的最小拷贝，只用到推任务通知需要的字段。
type wsChatMessage struct {
	Type      string `json:"type"`
	TaskID    string `json:"taskId,omitempty"`
	TaskName  string `json:"taskName,omitempty"`
	Status    string `json:"status,omitempty"`
	Title     string `json:"title,omitempty"`
	Level     string `json:"level,omitempty"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
}

func (c *TaskNotificationConsumer) pushWS(ctx context.Context, event taskNotificationEvent, title, content string) error {
	var msgType, level string
	switch event.Status {
	case consts.TaskStatusRunning:
		msgType, level = consts.WSTaskProgress, consts.TaskNotificationLevelInfo
	case consts.TaskStatusCompleted:
		msgType, level = consts.WSNotification, consts.TaskNotificationLevelSuccess
	case consts.TaskStatusFailed:
		msgType, level = consts.WSNotification, consts.TaskNotificationLevelError
	default:
		return nil
	}

	chatMsg := &wsChatMessage{
		Type:      msgType,
		TaskID:    fmt.Sprintf("%d", event.TaskID),
		TaskName:  event.TaskName,
		Status:    fmt.Sprintf("%d", event.Status),
		Title:     title,
		Level:     level,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}
	messageBytes, err := json.Marshal(chatMsg)
	if err != nil {
		return fmt.Errorf("任务通知消息序列化失败: %w", err)
	}

	resp, err := c.chatRPC.PushToUser(ctx, &chat.PushToUserRequest{UserId: event.UserID, PayloadJson: string(messageBytes)})
	if err != nil {
		logx.Errorf("任务通知回调 chat-rpc.PushToUser 失败: taskId=%d, userId=%d, err=%v", event.TaskID, event.UserID, err)
		return nil // 推送失败不影响主流程（尽力而为，语义与拆分前一致）
	}
	if !resp.Delivered {
		logx.Infof("任务通知WebSocket发送失败（用户可能不在线）: taskId=%d, userId=%d", event.TaskID, event.UserID)
	}
	return nil
}

func notificationText(event taskNotificationEvent) (title, content string, ok bool) {
	switch event.Status {
	case consts.TaskStatusRunning:
		return consts.TaskNotificationTitleRunning, fmt.Sprintf("任务「%s」正在执行中...", event.TaskName), true
	case consts.TaskStatusCompleted:
		return consts.TaskNotificationTitleCompleted, fmt.Sprintf("任务「%s」执行完成", event.TaskName), true
	case consts.TaskStatusFailed:
		return consts.TaskNotificationTitleFailed, fmt.Sprintf("任务「%s」执行失败：%s", event.TaskName, event.ErrorMessage), true
	default:
		return "", "", false
	}
}

// handleFailure 处理失败超过 maxDeliveryAttempts 次的消息：转存进 deadletter 流，
// XACK 原消息避免无限阻塞后续消息处理。见 17-async-eventing.md 第 2.4 节。
func (c *TaskNotificationConsumer) handleFailure(ctx context.Context, msgID string, event taskNotificationEvent, procErr error) {
	c.failuresMu.Lock()
	c.failures[msgID]++
	attempts := c.failures[msgID]
	c.failuresMu.Unlock()

	logx.Errorf("处理 task 通知事件失败(第 %d 次): msgID=%s, taskId=%d, error: %v", attempts, msgID, event.TaskID, procErr)

	if attempts < maxDeliveryAttempts {
		return // 不 ACK，留在 pending，等下次 XREADGROUP 重新投递
	}

	deadletterStream := StreamTaskNotification + streamDeadLetterSuffix
	payload, _ := json.Marshal(event)
	if _, err := c.redis.XAddCtx(ctx, deadletterStream, false, "*", []string{"payload", string(payload), "error", procErr.Error()}); err != nil {
		logx.Errorf("写入死信流失败: msgID=%s, error: %v", msgID, err)
	}

	c.failuresMu.Lock()
	delete(c.failures, msgID)
	c.failuresMu.Unlock()

	c.ack(ctx, msgID)
}

func (c *TaskNotificationConsumer) ack(ctx context.Context, msgID string) {
	if _, err := c.redis.XAckCtx(ctx, StreamTaskNotification, consumerGroup, msgID); err != nil {
		logx.Errorf("XACK 失败: msgID=%s, error: %v", msgID, err)
	}
}
