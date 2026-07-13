package middleware

import (
	"net/http"
	"strings"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/pkg/response"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// AuthMiddleware 校验 Access Token + 黑名单，并将用户信息写入 context。iam-rpc 拆分后
// 黑名单校验直连共享 Redis，不走 RPC（见 16-rpc-conventions.md 第 6 节：黑名单 100%
// 基于 Redis Exists，没有触碰任何 MySQL 表，iam-rpc 侧的 Logout/Refresh 往同一个共享
// Redis 实例写入黑名单 key，gateway 这里直接读，热路径零 RPC）。
type AuthMiddleware struct {
	redis     *redis.Redis
	jwtConfig config.JWTConf
}

func NewAuthMiddleware(cfg config.Config, rdb *redis.Redis) *AuthMiddleware {
	return &AuthMiddleware{redis: rdb, jwtConfig: cfg.JWT}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "未提供认证信息"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "无效的认证头"))
			return
		}
		token := parts[1]

		// 黑名单校验：直连共享 Redis，key 格式与 services/iam 内部的
		// TokenBlacklistRepository 保持一致（两边各自维护一份常量，不共享包，见
		// 16-rpc-conventions.md 第 6 节"直接复制不共享"）。
		blacklisted, err := m.redis.Exists(consts.RedisJWTBlacklistPrefix + token)
		if err != nil {
			response.ErrorCtx(r.Context(), w, errs.Wrap(errs.CodeInternalError, "检查令牌黑名单失败", err))
			return
		}
		if blacklisted {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "令牌已失效"))
			return
		}

		// 解析 Access Token
		claims, err := jwthelper.ParseToken(token, m.jwtConfig.AccessSecret)
		if err != nil || claims.IsRefresh {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "访问令牌无效或已过期"))
			return
		}

		ctxWithUser := jwthelper.WithAuthUser(r.Context(), jwthelper.AuthUser{
			UserID:   claims.UserID,
			Username: claims.Username,
		})

		next(w, r.WithContext(ctxWithUser))
	}
}
