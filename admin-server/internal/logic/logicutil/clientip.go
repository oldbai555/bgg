package logicutil

import (
	"net/http"
	"strings"

	"postapocgame/admin-server/internal/consts"
)

// ClientIP 取客户端 IP：X-Forwarded-For 首段 > X-Real-IP > RemoteAddr（去端口）。
func ClientIP(r *http.Request) string {
	if ip := r.Header.Get(consts.HeaderXForwardedFor); ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	if ip := r.Header.Get(consts.HeaderXRealIP); ip != "" {
		return ip
	}
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
