// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package metric

import (
	"net/http"

	"postapocgame/admin-server/internal/svc"
)

// MetricReportOptionsHandler 打点上报接口 CORS 预检处理
// 支持跨域请求，允许前端页面发起埋点上报
func MetricReportOptionsHandler(_ *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头，允许所有域名（可根据需要配置白名单）
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24小时
		w.WriteHeader(http.StatusNoContent)
	}
}
