// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video_collect

import (
	"net/http"

	"postapocgame/admin-server/internal/svc"
)

// VideoCollectOptionsHandler 采集接口 CORS 预检处理
func VideoCollectOptionsHandler(_ *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 这里仅处理预检请求并返回允许的 CORS 头
		w.Header().Set("Access-Control-Allow-Origin", "https://missav.ai")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
	}
}
