package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DSN string
	}
	// ChatRedis 不叫 Redis：zrpc.RpcServerConf 内嵌字段本身就有一个 Redis（用于可选的 gRPC
	// 鉴权），撞名会导致 go-zero 的 conf 解析绑定到错误的结构体，和 services/task、
	// services/sdk 的 TaskRedis/SdkRedis 同一个坑。这里的 Redis 不只是满足 goctl 生成 Model
	// 的 CachedConn 缓存节点要求，还真正用于 Redis Streams（stream:chat.user.created 消费者）
	// ——和 gateway 共享同一个 Redis 实例（缓存/锁/队列不拆分，见 16-rpc-conventions.md 第 6 节）。
	ChatRedis struct {
		Address  string
		Password string
	}
	// IamCallbackRpc 连到单体内嵌 pkg/iamcallback.IamCallback server 的 zrpc client 配置。
	// iam 域还没拆分成独立服务前的临时方案，见 pkg/iamcallback 包注释、
	// internal/rpcserver/iamcallback/server.go。
	IamCallbackRpc zrpc.RpcClientConf
}
