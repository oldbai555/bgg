// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package metric

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/monitoring/metric"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

type contextKey string

const (
	// IMPORTANT: use custom type for context keys to avoid collisions
	ctxKeyClientIP  contextKey = "client_ip"
	ctxKeyUserAgent contextKey = "user_agent"
)

func MetricReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MetricReportReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 提取客户端 IP 和 User-Agent，通过 context 传递到 logic
		clientIP := getClientIP(r)
		userAgent := r.UserAgent()
		ctx := context.WithValue(r.Context(), ctxKeyClientIP, clientIP)
		ctx = context.WithValue(ctx, ctxKeyUserAgent, userAgent)

		l := metric.NewMetricReportLogic(ctx, svcCtx)
		resp, err := l.MetricReport(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// getClientIP 获取客户端 IP 地址
func getClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	// 优先从 X-Forwarded-For 获取（代理场景）
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 其次从 X-Real-IP 获取
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 最后从 RemoteAddr 获取
	ip = r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
