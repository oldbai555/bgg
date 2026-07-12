// Package consumer 消费 IAM 侧（internal/domain/iam/user_service.go）发布的
// stream:chat.user.created 事件，触发新用户的 Chat onboarding（加入默认群 + 为存量用户
// 建私聊）。结构与 internal/consumer/task_notification_consumer.go 完全同构（同一套
// XGROUP/XREADGROUP/XACK/死信处理骨架，见 17-async-eventing.md 第 2.4/2.5 节的幂等/死信
// 要求），这是本仓库第一次真正在生产路径上使用 Streams（task-rpc 拆分阶段实际用的是
// TaskCallback 同步 RPC，不是 Streams；stream:task.notification 消费者虽然更早写好，但
// 和这里是同一批次落地的两个独立消费者，互不依赖）。
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

	chatdomain "postapocgame/admin-server/services/chat/internal/domain/chat"
)

const (
	// StreamChatUserCreated 必须和 internal/domain/iam/user_service.go 里的同名常量保持
	// 一致，两边各自维护一份（16-rpc-conventions.md 第 6 节"直接复制不共享"）。
	StreamChatUserCreated  = "stream:chat.user.created"
	streamDeadLetterSuffix = ".deadletter"
	consumerGroup          = "chat-rpc-init"
	maxDeliveryAttempts    = 3
)

// chatUserCreatedEvent 与 IAM 侧 ChatUserCreatedEvent 字段一一对应。
type chatUserCreatedEvent struct {
	UserID    uint64 `json:"userId"`
	CreatedAt int64  `json:"createdAt"`
}

// ChatUserCreatedConsumer 消费 stream:chat.user.created。
type ChatUserCreatedConsumer struct {
	redis      *redis.Redis
	onboarding *chatdomain.ChatOnboardingService
	stopChan   chan struct{}
	wg         sync.WaitGroup

	consumerName string

	failuresMu sync.Mutex
	failures   map[string]int
}

func NewChatUserCreatedConsumer(rds *redis.Redis, onboarding *chatdomain.ChatOnboardingService) *ChatUserCreatedConsumer {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "unknown-host"
	}
	return &ChatUserCreatedConsumer{
		redis:        rds,
		onboarding:   onboarding,
		stopChan:     make(chan struct{}),
		failures:     make(map[string]int),
		consumerName: fmt.Sprintf("%s-%s-%d", consumerGroup, hostname, os.Getpid()),
	}
}

// Start 启动消费者组（幂等，BUSYGROUP 错误直接忽略）并起后台 goroutine 开始消费。
func (c *ChatUserCreatedConsumer) Start() {
	if _, err := c.redis.XGroupCreateMkStreamCtx(context.Background(), StreamChatUserCreated, consumerGroup, "$"); err != nil {
		logx.Infof("XGROUP CREATE %s %s: %v（已存在则忽略）", StreamChatUserCreated, consumerGroup, err)
	}

	c.wg.Add(1)
	go c.run()
	logx.Infof("chat onboarding 消费者已启动: stream=%s, group=%s", StreamChatUserCreated, consumerGroup)
}

// Stop 停止消费者。
func (c *ChatUserCreatedConsumer) Stop() {
	close(c.stopChan)
	c.wg.Wait()
	logx.Infof("chat onboarding 消费者已停止")
}

func (c *ChatUserCreatedConsumer) run() {
	defer c.wg.Done()

	node, err := redis.CreateBlockingNode(c.redis)
	if err != nil {
		logx.Errorf("chat onboarding 消费者创建 blocking node 失败: %v", err)
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
			StreamChatUserCreated, ">")
		if err != nil {
			if err != redis.Nil {
				logx.Errorf("XREADGROUP %s 失败: %v", StreamChatUserCreated, err)
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

func (c *ChatUserCreatedConsumer) handleMessage(ctx context.Context, msgID string, values map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("chat onboarding 消费发生 panic: msgID=%s, error: %v", msgID, r)
		}
	}()

	payload, _ := values["payload"].(string)
	var event chatUserCreatedEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		logx.Errorf("解析 chat.user.created 事件失败，直接丢弃: msgID=%s, error: %v", msgID, err)
		c.ack(ctx, msgID)
		return
	}

	// InitNewUser 内部的 joinDefaultGroup/createPrivateChatsForExistingUsers 本身就是
	// "先查是否已存在再插入"（见 onboarding.go 注释），天然幂等，同一条消息被重复投递
	// 不会产生重复的群成员/私聊记录，符合 17-async-eventing.md 第 2.4 节的幂等要求。
	if err := c.onboarding.InitNewUser(ctx, event.UserID); err != nil {
		c.handleFailure(ctx, msgID, event, err)
		return
	}

	c.failuresMu.Lock()
	delete(c.failures, msgID)
	c.failuresMu.Unlock()
	c.ack(ctx, msgID)
}

// handleFailure 处理失败超过 maxDeliveryAttempts 次的消息：转存进 deadletter 流，
// XACK 原消息避免无限阻塞后续消息处理。见 17-async-eventing.md 第 2.5 节。
func (c *ChatUserCreatedConsumer) handleFailure(ctx context.Context, msgID string, event chatUserCreatedEvent, procErr error) {
	c.failuresMu.Lock()
	c.failures[msgID]++
	attempts := c.failures[msgID]
	c.failuresMu.Unlock()

	logx.Errorf("处理 chat.user.created 事件失败(第 %d 次): msgID=%s, userId=%d, error: %v", attempts, msgID, event.UserID, procErr)

	if attempts < maxDeliveryAttempts {
		return // 不 ACK，留在 pending，等下次 XREADGROUP 重新投递
	}

	deadletterStream := StreamChatUserCreated + streamDeadLetterSuffix
	payload, _ := json.Marshal(event)
	if _, err := c.redis.XAddCtx(ctx, deadletterStream, false, "*", []string{"payload", string(payload), "error", procErr.Error()}); err != nil {
		logx.Errorf("写入死信流失败: msgID=%s, error: %v", msgID, err)
	}

	c.failuresMu.Lock()
	delete(c.failures, msgID)
	c.failuresMu.Unlock()

	c.ack(ctx, msgID)
}

func (c *ChatUserCreatedConsumer) ack(ctx context.Context, msgID string) {
	if _, err := c.redis.XAckCtx(ctx, StreamChatUserCreated, consumerGroup, msgID); err != nil {
		logx.Errorf("XACK 失败: msgID=%s, error: %v", msgID, err)
	}
}
