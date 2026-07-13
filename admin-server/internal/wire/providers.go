package wire

import (
	"os"
	"path/filepath"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/middleware"
	"postapocgame/admin-server/internal/redisconn"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/pkg/iamcallback"
	"postapocgame/admin-server/services/chat/chatclient"
	"postapocgame/admin-server/services/content/contentclient"
	"postapocgame/admin-server/services/iam/iamclient"
	"postapocgame/admin-server/services/sdk/sdkclient"
	"postapocgame/admin-server/services/task/taskclient"

	"github.com/google/wire"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

// ProviderSet 组合根依赖注入集合。iam+system+monitoring+misc 四个域拆分成 iam-rpc 后，
// provideRepository/provideDomain 两个 provider 整个删除，换成 provideRedis（gateway
// 唯一还需要直连的基础设施：token 黑名单 + 限流滑动窗口）。PermissionMiddleware 的构造
// 函数现在直接吃 iamclient.Iam（已经是独立 Wire 节点），不再需要 providePermissionMiddleware
// 这层适配函数。
var ProviderSet = wire.NewSet(
	provideRedis,
	provideIamRPC,
	provideIamCallbackRPC,
	provideTaskRPC,
	provideSdkRPC,
	provideChatRPC,
	provideContentRPC,

	middleware.NewAuthMiddleware,
	middleware.NewApiEnabledMiddleware,
	middleware.NewPermissionMiddleware,
	middleware.NewOperationLogMiddleware,
	middleware.NewPublicOperationLogMiddleware,
	middleware.NewRateLimitMiddleware,
	middleware.NewPerformanceMiddleware,
	middleware.NewCorsMiddleware,
	middleware.NewSDKAuthMiddleware,
	middleware.NewSDKRateLimitMiddleware,
	middleware.NewSDKCallLogMiddleware,
	provideMiddlewareBundle,

	provideServiceContext,
)

// provideRedis 是 gateway 拆分完 iam-rpc 后唯一还需要的直连基础设施：token 黑名单
// （AuthMiddleware）+ 限流滑动窗口（RateLimitMiddleware/SDKRateLimitMiddleware）。
func provideRedis(c config.Config) (*redis.Redis, error) {
	if err := initUploadDir(); err != nil {
		return nil, err
	}
	return redisconn.New(c.Redis)
}

func initUploadDir() error {
	if err := os.MkdirAll(consts.UploadDir, 0o755); err != nil {
		logx.Errorf("创建 uploads 目录失败: %v", err)
		return err
	}
	absPath, _ := filepath.Abs(consts.UploadDir)
	logx.Infof("uploads 目录已初始化: %s", absPath)
	return nil
}

// provideIamRPC 连到 iam-rpc（services/iam/）。iam+system+monitoring+misc 四个域已经
// 拆分成独立服务，gateway 侧不再直接持有 Repository/Domain，改成一个 zrpc client。
func provideIamRPC(c config.Config) iamclient.Iam {
	return iamclient.NewIam(zrpc.MustNewClient(c.IamRpc))
}

// provideIamCallbackRPC 连到 iam-rpc 同一进程内注册的 pkg/iamcallback.IamCallback 服务
// （RBAC 变更审计日志写入，pkg/audit.RecordAuditLog 用它）。
func provideIamCallbackRPC(c config.Config) (iamcallback.Client, error) {
	return iamcallback.NewClient(c.IamCallbackRpc)
}

// provideTaskRPC 连到 task-rpc（services/task/）。task 域已经拆分成独立服务，gateway
// 侧不再直接持有 TaskExecutors/TaskScheduler，改成一个 zrpc client。
func provideTaskRPC(c config.Config) taskclient.Task {
	return taskclient.NewTask(zrpc.MustNewClient(c.TaskRPCConf))
}

// provideSdkRPC 连到 sdk-rpc（services/sdk/）。sdk 域已经拆分成独立服务，gateway 侧
// 不再直接持有 Domain.SDK，改成一个 zrpc client。
func provideSdkRPC(c config.Config) sdkclient.Sdk {
	return sdkclient.NewSdk(zrpc.MustNewClient(c.SdkRPCConf))
}

// provideChatRPC 连到 chat-rpc（services/chat/）。chat 域已经拆分成独立服务，gateway 侧
// 不再直接持有 Domain.Chat/ChatHub，改成一个 zrpc client。
func provideChatRPC(c config.Config) chatclient.Chat {
	return chatclient.NewChat(zrpc.MustNewClient(c.ChatRPCConf))
}

// provideContentRPC 连到 content-rpc（services/content/）。blog+video 域已经拆分成独立
// 服务，gateway 侧不再直接持有 Domain.Blog/Domain.Video，改成一个 zrpc client。
func provideContentRPC(c config.Config) contentclient.Content {
	return contentclient.NewContent(zrpc.MustNewClient(c.ContentRPCConf))
}

// provideMiddlewareBundle 是唯一还需要手写的 assembler：MiddlewareBundle 的 11 个字段
// 类型全部是 rest.Middleware（同一个具名类型），Wire 无法自动区分该把哪个 provider
// 的结果填进哪个字段，所以仍然需要一个手写函数。
func provideMiddlewareBundle(
	auth *middleware.AuthMiddleware,
	apiEnabled *middleware.ApiEnabledMiddleware,
	permission *middleware.PermissionMiddleware,
	operationLog *middleware.OperationLogMiddleware,
	publicOperationLog *middleware.PublicOperationLogMiddleware,
	rateLimit *middleware.RateLimitMiddleware,
	performance *middleware.PerformanceMiddleware,
	cors *middleware.CorsMiddleware,
	sdkAuth *middleware.SDKAuthMiddleware,
	sdkRateLimit *middleware.SDKRateLimitMiddleware,
	sdkCallLog *middleware.SDKCallLogMiddleware,
) *MiddlewareBundle {
	return &MiddlewareBundle{
		Auth:               auth.Handle,
		ApiEnabled:         apiEnabled.Handle,
		Permission:         permission.Handle,
		OperationLog:       operationLog.Handle,
		PublicOperationLog: publicOperationLog.Handle,
		RateLimit:          rateLimit.Handle,
		Performance:        performance.Handle,
		Cors:               cors.Handle,
		SDKAuth:            sdkAuth.Handle,
		SDKRateLimit:       sdkRateLimit.Handle,
		SDKCallLog:         sdkCallLog.Handle,
	}
}

func provideServiceContext(
	c config.Config,
	rdb *redis.Redis,
	iamRPC iamclient.Iam,
	iamCallbackRPC iamcallback.Client,
	taskRPC taskclient.Task,
	sdkRPC sdkclient.Sdk,
	chatRPC chatclient.Chat,
	contentRPC contentclient.Content,
	mw *MiddlewareBundle,
) (*svc.ServiceContext, func()) {
	svcCtx := &svc.ServiceContext{
		Config:                       c,
		Redis:                        rdb,
		IamRPC:                       iamRPC,
		IamCallbackRPC:               iamCallbackRPC,
		TaskRPC:                      taskRPC,
		SdkRPC:                       sdkRPC,
		ChatRPC:                      chatRPC,
		ContentRPC:                   contentRPC,
		AuthMiddleware:               mw.Auth,
		ApiEnabledMiddleware:         mw.ApiEnabled,
		PermissionMiddleware:         mw.Permission,
		OperationLogMiddleware:       mw.OperationLog,
		PublicOperationLogMiddleware: mw.PublicOperationLog,
		RateLimitMiddleware:          mw.RateLimit,
		PerformanceMiddleware:        mw.Performance,
		CorsMiddleware:               mw.Cors,
		SDKAuthMiddleware:            mw.SDKAuth,
		SDKRateLimitMiddleware:       mw.SDKRateLimit,
		SDKCallLogMiddleware:         mw.SDKCallLog,
	}

	cleanup := func() {}
	return svcCtx, cleanup
}
