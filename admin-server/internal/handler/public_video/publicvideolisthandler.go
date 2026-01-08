// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public_video

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/public_video"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func PublicVideoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublicVideoListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public_video.NewPublicVideoListLogic(r.Context(), svcCtx)
		resp, err := l.PublicVideoList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
