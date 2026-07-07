// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package video_collect

import (
	"net/http"

	"postapocgame/admin-server/internal/logic/video/video_collect"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func VideoCollectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头，支持跨域请求
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		var req types.VideoCollectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := video_collect.NewVideoCollectLogic(r.Context(), svcCtx)
		resp, err := l.VideoCollect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
