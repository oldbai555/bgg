// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package my

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/fitness/my"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func FitnessMySuggestionDecideHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FitnessMySuggestionDecideReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := my.NewFitnessMySuggestionDecideLogic(r.Context(), svcCtx)
		err := l.FitnessMySuggestionDecide(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
