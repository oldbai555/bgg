// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package m3u8

import (
	"net/http"

	"postapocgame/admin-server/internal/logic/m3u8"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func M3u8ProxyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 注意：OPTIONS 预检请求已由 nginx 在 /gateway/ location 中处理并返回 204
		// CORS 响应头也由 nginx 统一设置（使用 always 参数，确保错误响应也有 CORS 头）
		// 这里不再重复处理，简化代码

		var req types.M3u8ProxyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := m3u8.NewM3u8ProxyLogic(r.Context(), svcCtx)
		err := l.M3u8Proxy(w, r, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		// m3u8代理直接写入响应流，不需要返回 JSON
	}
}
