package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DSN string
	}
	// IamRedis 不叫 Redis：zrpc.RpcServerConf 内嵌字段本身就有一个 Redis（用于可选的 gRPC
	// 鉴权），撞名会导致 go-zero 的 conf 解析绑定到错误的结构体，和 task/sdk/chat/content 的
	// TaskRedis/SdkRedis/ChatRedis/ContentRedis 同一个坑。iam-rpc 是这个共享 Redis 实例
	// 真正的读写方（JWT 黑名单登记、BusinessCache 权限缓存、Token 刷新等），与 gateway/
	// 其余四个服务共享同一个物理实例（缓存/锁/队列不拆分，见 16-rpc-conventions.md 第 6 节）。
	IamRedis struct {
		Address  string
		Password string
	}
	JWT struct {
		AccessSecret  string
		RefreshSecret string
		AccessExpire  int64
		RefreshExpire int64
		Issuer        string
	}
	Bcrypt struct {
		Cost int
	}
	// SdkRpc 连到 sdk-rpc 的 zrpc client 配置：原单体内嵌 TaskCallback server 的
	// fetchSdkCallLog 分支回调 sdk-rpc.SdkCallLogExport，这里原样保留（见
	// internal/rpcserver/taskcallback/server.go 搬迁前的实现）。
	SdkRpc zrpc.RpcClientConf
	// ChatRpc 连到 chat-rpc 的 zrpc client 配置：task 通知消费者（原单体
	// internal/consumer/task_notification_consumer.go，现原样搬到
	// services/iam/internal/consumer/）推送 WS 通知时回调 chat-rpc.PushToUser。
	ChatRpc zrpc.RpcClientConf
}
