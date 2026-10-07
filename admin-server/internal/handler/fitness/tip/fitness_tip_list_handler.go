// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tip

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/fitness/tip"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func FitnessTipListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FitnessTipListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tip.NewFitnessTipListLogic(r.Context(), svcCtx)
		resp, err := l.FitnessTipList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
