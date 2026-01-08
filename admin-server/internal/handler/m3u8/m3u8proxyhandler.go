// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package m3u8

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/m3u8"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func M3u8ProxyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
