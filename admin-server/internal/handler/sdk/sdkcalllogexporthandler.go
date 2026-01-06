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

func SdkCallLogExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SdkCallLogExportReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sdk.NewSdkCallLogExportLogic(r.Context(), svcCtx)
		err := l.SdkCallLogExport(w, r, &req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// 导出功能直接写入响应流，不需要返回 JSON
	}
}
