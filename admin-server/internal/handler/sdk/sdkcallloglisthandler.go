// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/sdk"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func SdkCallLogListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SdkCallLogListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sdk.NewSdkCallLogListLogic(r.Context(), svcCtx)
		resp, err := l.SdkCallLogList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
