// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/response"
)

type sdkCtxKey string

const (
	ctxKeySdkKeyID       sdkCtxKey = "sdkKeyId"
	ctxKeySdkApiKey      sdkCtxKey = "sdkApiKey"
	ctxKeySdkInterfaceID sdkCtxKey = "sdkInterfaceId"
	ctxKeySdkApiCode     sdkCtxKey = "sdkApiCode"
)

type SDKAuthMiddleware struct {
	svcCtx *svc.ServiceContext
}

func NewSDKAuthMiddleware(svcCtx *svc.ServiceContext) *SDKAuthMiddleware {
	return &SDKAuthMiddleware{svcCtx: svcCtx}
}

func (m *SDKAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		apiSecret := r.Header.Get("X-API-Secret")
		if strings.TrimSpace(apiKey) == "" || strings.TrimSpace(apiSecret) == "" {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "缺少 API Key 或 Secret"))
			return
		}

		sdkRepo := repository.NewSdkRepository(m.svcCtx.Repository)
		sdkKey, err := sdkRepo.FindKeyByApiKey(r.Context(), apiKey)
		if err != nil {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "无效的 API Key"))
			return
		}
		if sdkKey.ApiSecret != apiSecret {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "Secret 不匹配"))
			return
		}
		if sdkKey.Status != 1 {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "API Key 已被禁用"))
			return
		}
		if sdkKey.ExpireAt > 0 && time.Now().Unix() > sdkKey.ExpireAt {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "API Key 已过期"))
			return
		}

		clientIP := clientIPFromRequest(r)
		if !sdkRepo.IsIPAllowed(sdkKey.IpWhitelist, clientIP) {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeForbidden, "IP 不在白名单"))
			return
		}

		apiCode := sdkRepo.BuildInterfaceCode(r.Method, r.URL.Path)
		sdkInterface, err := sdkRepo.FindInterfaceByCode(r.Context(), apiCode)
		if err != nil || sdkInterface == nil || sdkInterface.Status != 1 {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeForbidden, "接口未开通或已禁用"))
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeySdkKeyID, sdkKey.Id)
		ctx = context.WithValue(ctx, ctxKeySdkApiKey, sdkKey.ApiKey)
		ctx = context.WithValue(ctx, ctxKeySdkInterfaceID, sdkInterface.Id)
		ctx = context.WithValue(ctx, ctxKeySdkApiCode, sdkInterface.ApiCode)
		next(w, r.WithContext(ctx))
	}
}

func clientIPFromRequest(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
