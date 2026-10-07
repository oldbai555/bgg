// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package template

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/fitness/template"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func FitnessTemplateUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FitnessTemplateUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := template.NewFitnessTemplateUpdateLogic(r.Context(), svcCtx)
		err := l.FitnessTemplateUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
