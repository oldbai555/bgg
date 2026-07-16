// Package redisconn 是 gateway 拆分完 iam-rpc 后剩下的唯一直连基础设施：共享 Redis
// （token 黑名单、限流滑动窗口，见 16-rpc-conventions.md 第 6 节"Redis 保持全服务共享"）。
// 不再需要 internal/repository.Repository 那一整套聚合 Model 的大句柄。
package redisconn

import (
	"github.com/zeromicro/go-zero/core/stores/redis"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/pkg/errs"
)

func New(cfg config.RedisConf) (*redis.Redis, error) {
	if cfg.Address == "" {
		return nil, errs.New(errs.CodeBadRequest, "redis address is empty")
	}
	return redis.NewRedis(redis.RedisConf{
		Host: cfg.Address,
		Pass: cfg.Password,
		Type: "node",
	})
}
