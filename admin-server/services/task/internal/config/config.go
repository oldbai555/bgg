package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DSN string
	}
	// TaskRedis 不叫 Redis：zrpc.RpcServerConf 内嵌字段本身就有一个 Redis
	// （redis.RedisKeyConf，用于可选的 gRPC 鉴权），撞名会导致 go-zero 的 conf 解析
	// 绑定到错误的结构体，YAML 里明明填了 Redis.Address 却报 "redis.Host is not set"。
	TaskRedis struct {
		Address  string
		Password string
	}
	// TaskCallbackRpc 连到承载导出数据/admin_file 的服务（当前阶段：单体内嵌的
	// pkg/taskcallback server）。见 pkg/taskcallback/taskcallback.proto。
	TaskCallbackRpc zrpc.RpcClientConf
	// RecentTaskLimit 是 TaskRecent 接口在请求未显式传 limit 时使用的默认条数，取代原来
	// 网关侧"缓存优先、字典兜底"的查询链路，见 15-service-boundaries.md 第 5 节末尾建议。
	RecentTaskLimit int64 `json:",default=10"`
}
