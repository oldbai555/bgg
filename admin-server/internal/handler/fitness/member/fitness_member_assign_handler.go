// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package member

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/fitness/member"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func FitnessMemberAssignHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FitnessMemberAssignReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := member.NewFitnessMemberAssignLogic(r.Context(), svcCtx)
		resp, err := l.FitnessMemberAssign(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
