package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DSN string
	}
	// SdkRedis 不叫 Redis：zrpc.RpcServerConf 内嵌字段本身就有一个 Redis
	// （redis.RedisKeyConf，用于可选的 gRPC 鉴权），撞名会导致 go-zero 的 conf 解析
	// 绑定到错误的结构体，和 services/task/internal/config/config.go 的 TaskRedis
	// 同一个坑。sdk-rpc 业务本身不用 Redis，这里纯粹是满足 goctl 生成的 Model 内部
	// CachedConn 强制要求非空缓存节点（cache.New 对空 CacheConf 会 log.Fatal），
	// 与 gateway 共享同一个 Redis 实例（缓存/锁/队列不拆分，见 16-rpc-conventions.md 第 6 节）。
	SdkRedis struct {
		Address  string
		Password string
	}
	// RateLimitDefault 取代原来读字典 sdk_rate_limit_default（物理属于 iam 域）的做法，
	// 见 services/sdk/internal/repository/sdk/sdk_repository.go 的 GetDefaultRateLimit
	// 注释、18-service-extraction-runbook.md 2.1/2.2 节。字典种子数据的默认值是 60。
	RateLimitDefault int64 `json:",default=60"`
}
