// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/response"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
)

type SDKRateLimitMiddleware struct {
	svcCtx *svc.ServiceContext
}

func NewSDKRateLimitMiddleware(svcCtx *svc.ServiceContext) *SDKRateLimitMiddleware {
	return &SDKRateLimitMiddleware{svcCtx: svcCtx}
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

		sdkRepo := sdkrepo.NewSdkRepository(m.svcCtx.Repository)
		iface, err := sdkRepo.FindInterfaceByCode(ctx, apiCode)
		if err != nil || iface == nil {
			response.ErrorCtx(ctx, w, errs.New(errs.CodeForbidden, "接口不存在"))
			return
		}
		binding, _ := sdkRepo.FindKeyApiBinding(ctx, sdkKeyId, sdkInterfaceId)

		limit := sdkRepo.GetDefaultRateLimit(ctx, iface.RateLimitDefault)
		if binding != nil && binding.CustomRateLimit > 0 {
			limit = binding.CustomRateLimit
		}
		if limit <= 0 {
			limit = 60
		}

		clientIP := clientIPFromRequest(r)
		now := time.Now()
		redis := m.svcCtx.Repository.Redis

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
				response.ErrorCtx(ctx, w, errs.New(errs.CodeForbidden, "请求过于频繁"))
				return
			}
		}

		next(w, r)
	}
}
