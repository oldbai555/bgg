package svc

import (
	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/repository/registry"
	"postapocgame/admin-server/services/chat/chatclient"
	"postapocgame/admin-server/services/content/contentclient"
	"postapocgame/admin-server/services/sdk/sdkclient"
	"postapocgame/admin-server/services/task/taskclient"

	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config     config.Config
	Repository *repository.Repository
	Domain     *registry.Domain
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
	ContentRPC                   contentclient.Content
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
