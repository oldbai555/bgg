// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/video"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func VideoProxyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VideoProxyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video.NewVideoProxyLogic(r.Context(), svcCtx)
		err := l.VideoProxy(w, r, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		// 视频代理直接写入响应流，不需要返回 JSON
	}
}
