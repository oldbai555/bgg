// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/iam/auth"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func LoginFeishuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginFeishuReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth.NewLoginFeishuLogic(r.Context(), svcCtx)
		resp, err := l.LoginFeishu(&req, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
