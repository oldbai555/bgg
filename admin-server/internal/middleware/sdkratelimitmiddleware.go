// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/response"
	"postapocgame/admin-server/services/sdk/sdkclient"
)

// SDKRateLimitMiddleware 继续留在 gateway：限流计数用的 Redis 滑动窗口是全服务共享的
// 基础设施，不属于 sdk 域数据，留在这里；只有"这个 Key/接口组合的有效限流上限是多少"这一
// 步查询改成调 sdk-rpc 的 GetEffectiveRateLimit（sdk_interface.rate_limit_default +
// sdk_key_api.custom_rate_limit 覆盖，都是 sdk 域自己的表）。见 18-service-extraction-
// runbook.md 2.2 节。
type SDKRateLimitMiddleware struct {
	repo   *repository.Repository
	sdkRPC sdkclient.Sdk
}

func NewSDKRateLimitMiddleware(repo *repository.Repository, sdkRPC sdkclient.Sdk) *SDKRateLimitMiddleware {
	return &SDKRateLimitMiddleware{repo: repo, sdkRPC: sdkRPC}
}

func (m *SDKRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sdkKeyId, _ := ctx.Value(ctxKeySdkKeyID).(uint64)
		apiCode, _ := ctx.Value(ctxKeySdkApiCode).(string)
		sdkInterfaceId, _ := ctx.Value(ctxKeySdkInterfaceID).(uint64)
		apiKey, _ := ctx.Value(ctxKeySdkApiKey).(string)
		if sdkKeyId == 0 || sdkInterfaceId == 0 || apiCode == "" || apiKey == "" {
			response.ErrorCtx(ctx, w, errs.New(errs.CodeForbidden, "SDK 鉴权信息缺失"))
			return
		}

		rpcResp, err := m.sdkRPC.GetEffectiveRateLimit(ctx, &sdkclient.GetEffectiveRateLimitRequest{
			SdkKeyId:       sdkKeyId,
			SdkInterfaceId: sdkInterfaceId,
		})
		if err != nil {
			response.ErrorCtx(ctx, w, errs.WrapGRPCError("获取限流配置失败", err))
			return
		}
		limit := rpcResp.Limit
		if limit <= 0 {
			limit = 60
		}

		clientIP := clientIPFromRequest(r)
		now := time.Now()
		redis := m.repo.Redis

		keys := []string{
			fmt.Sprintf("sdk:rl:key:%s:%s:%d", apiKey, apiCode, now.Unix()/60),
			fmt.Sprintf("sdk:rl:ip:%s:%s:%d", clientIP, apiCode, now.Unix()/60),
		}
		for _, k := range keys {
			cnt, err := redis.IncrCtx(ctx, k)
			if err != nil {
				continue
			}
			if cnt == 1 {
				_ = redis.ExpireCtx(ctx, k, 65)
			}
			if cnt > int64(limit) {
				// 与 Admin 侧 RateLimitMiddleware 保持一致：先手动写 429（response.ErrorCtx
				// 对业务错误统一写 400，net/http 对同一响应重复 WriteHeader 只认第一次）。
				w.WriteHeader(http.StatusTooManyRequests)
				response.ErrorCtx(ctx, w, errs.New(errs.CodeTooManyRequests, "请求过于频繁"))
				return
			}
		}

		next(w, r)
	}
}
