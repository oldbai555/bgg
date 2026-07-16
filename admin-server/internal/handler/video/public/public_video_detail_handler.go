// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/video/public"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func PublicVideoDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublicVideoDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public.NewPublicVideoDetailLogic(r.Context(), svcCtx)
		resp, err := l.PublicVideoDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
