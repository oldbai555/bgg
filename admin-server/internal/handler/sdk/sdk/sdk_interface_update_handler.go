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

func SdkInterfaceUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SdkInterfaceUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sdk.NewSdkInterfaceUpdateLogic(r.Context(), svcCtx)
		err := l.SdkInterfaceUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
