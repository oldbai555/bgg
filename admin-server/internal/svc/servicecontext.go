package svc

import (
	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/ollama"
	"postapocgame/admin-server/pkg/iamcallback"
	"postapocgame/admin-server/services/chat/chatclient"
	"postapocgame/admin-server/services/content/contentclient"
	"postapocgame/admin-server/services/iam/iamclient"
	"postapocgame/admin-server/services/sdk/sdkclient"
	"postapocgame/admin-server/services/task/taskclient"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config config.Config
	// Redis 是全服务共享的 Redis 客户端（token 黑名单、限流滑动窗口），iam+system+
	// monitoring+misc 域拆分成 iam-rpc 后，gateway 不再直连任何 MySQL，Repository/Domain
	// 两个聚合根字段整个删除。
	Redis *redis.Redis
	// IamRPC 是 iam-rpc（services/iam/）的 zrpc client，取代了原来直接持有的
	// Repository/Domain——iam+system+monitoring+misc 四个域已经拆分成独立服务，gateway
	// 侧全部相关 logic 和 5 个中间件（Auth/Permission/ApiEnabled/OperationLog/
	// Performance）都通过这个 client 调用。
	IamRPC iamclient.Iam
	// IamCallbackRPC 连到 iam-rpc 同一进程内注册的 pkg/iamcallback.IamCallback 服务，
	// 目前唯一的调用方是 pkg/audit.RecordAuditLog（RBAC 变更审计日志）。
	IamCallbackRPC iamcallback.Client
	// TaskRPC 是 task-rpc（services/task/）的 zrpc client，取代了原来直接持有的
	// TaskExecutors/TaskScheduler——task 域已经拆分成独立服务，gateway 侧只剩薄胶水。
	TaskRPC taskclient.Task
	// SdkRPC 是 sdk-rpc（services/sdk/）的 zrpc client，取代了原来直接持有的
	// Domain.SDK——sdk 域已经拆分成独立服务，11 个 SdkApiKey/SdkInterface/SdkCallLog
	// logic 和 SDKAuthMiddleware/SDKRateLimitMiddleware/SDKCallLogMiddleware 三个
	// 中间件都通过这个 client 调用。
	SdkRPC sdkclient.Sdk
	// ChatRPC 是 chat-rpc（services/chat/）的 zrpc client，取代了原来直接持有的
	// Domain.Chat/ChatHub——chat 域已经拆分成独立服务，11 个 Chat*/ChatGroup*/ChatMessage*
	// logic 和 WS↔gRPC 桥接 handler（internal/handler/chat/chatwshandler.go）都通过这个
	// client 调用。
	ChatRPC chatclient.Chat
	// ContentRPC 是 content-rpc（services/content/）的 zrpc client，取代了原来直接持有的
	// Domain.Blog/Domain.Video——blog+video 域已经拆分成独立服务，34 个 Blog*/PublicBlog*
	// logic + 6 个 Video*/PublicVideo* logic 都通过这个 client 调用（M3u8Proxy 纯 HTTP 代理
	// + VideoCollectOptions CORS 预检占位，不访问域数据，继续留在 gateway 不接入）。
	ContentRPC contentclient.Content
	// OllamaClient 连本机 Ollama REST API，供 ai/knowledge_qa 的 Reindex/Ask 两个 Logic
	// 做 embedding + 生成，详见 docs/ai-knowledge-qa-spec.md。
	OllamaClient                 *ollama.Client
	AuthMiddleware               rest.Middleware
	ApiEnabledMiddleware         rest.Middleware
	PermissionMiddleware         rest.Middleware
	OperationLogMiddleware       rest.Middleware
	PublicOperationLogMiddleware rest.Middleware
	RateLimitMiddleware          rest.Middleware
	PerformanceMiddleware        rest.Middleware
	CorsMiddleware               rest.Middleware
	SDKAuthMiddleware            rest.Middleware
	SDKRateLimitMiddleware       rest.Middleware
	SDKCallLogMiddleware         rest.Middleware
}
