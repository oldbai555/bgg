// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/response"
	"postapocgame/admin-server/services/sdk/sdkclient"
)

type sdkCtxKey string

const (
	ctxKeySdkKeyID       sdkCtxKey = "sdkKeyId"
	ctxKeySdkApiKey      sdkCtxKey = "sdkApiKey"
	ctxKeySdkInterfaceID sdkCtxKey = "sdkInterfaceId"
	ctxKeySdkApiCode     sdkCtxKey = "sdkApiCode"
)

// SDKAuthMiddleware 继续留在 gateway（HTTP 请求最早触达的地方），内部实现从直连
// Repository 改成调 sdk-rpc 的 VerifyApiKey，sdk 域的鉴权判断逻辑本身已经原样搬进
// services/sdk/internal/logic/verifyapikeylogic.go，见 18-service-extraction-runbook.md
// 2.2 节。
type SDKAuthMiddleware struct {
	sdkRPC sdkclient.Sdk
}

func NewSDKAuthMiddleware(sdkRPC sdkclient.Sdk) *SDKAuthMiddleware {
	return &SDKAuthMiddleware{sdkRPC: sdkRPC}
}

func (m *SDKAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		apiSecret := r.Header.Get("X-API-Secret")

		resp, err := m.sdkRPC.VerifyApiKey(r.Context(), &sdkclient.VerifyApiKeyRequest{
			ApiKey:    apiKey,
			ApiSecret: apiSecret,
			Method:    r.Method,
			Path:      r.URL.Path,
			ClientIp:  clientIPFromRequest(r),
		})
		if err != nil {
			response.ErrorCtx(r.Context(), w, errs.WrapGRPCError("SDK 鉴权失败", err))
			return
		}
		if !resp.Valid {
			response.ErrorCtx(r.Context(), w, errs.New(int(resp.Code), resp.Message))
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeySdkKeyID, resp.SdkKeyId)
		ctx = context.WithValue(ctx, ctxKeySdkApiKey, resp.ApiKey)
		ctx = context.WithValue(ctx, ctxKeySdkInterfaceID, resp.SdkInterfaceId)
		ctx = context.WithValue(ctx, ctxKeySdkApiCode, resp.ApiCode)
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
