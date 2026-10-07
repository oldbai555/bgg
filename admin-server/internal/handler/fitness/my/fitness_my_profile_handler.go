// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package my

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/fitness/my"
	"postapocgame/admin-server/internal/svc"
)

func FitnessMyProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := my.NewFitnessMyProfileLogic(r.Context(), svcCtx)
		resp, err := l.FitnessMyProfile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
