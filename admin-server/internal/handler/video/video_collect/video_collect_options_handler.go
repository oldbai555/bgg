// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package video_collect

import (
	"net/http"

	"postapocgame/admin-server/internal/svc"
)

// VideoCollectOptionsHandler 采集接口 CORS 预检处理
// 支持跨域请求，允许第三方页面发起采集请求
func VideoCollectOptionsHandler(_ *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头，允许所有域名（可根据需要配置白名单）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24小时
		w.WriteHeader(http.StatusNoContent)
	}
}
