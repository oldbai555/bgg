// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/sdk/sdk"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func SdkInterfaceDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SdkInterfaceDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sdk.NewSdkInterfaceDeleteLogic(r.Context(), svcCtx)
		err := l.SdkInterfaceDelete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
