// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/sdk/public"
	"postapocgame/admin-server/internal/svc"
)

func SdkFileUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewSdkFileUploadLogic(r.Context(), svcCtx)
		resp, err := l.SdkFileUpload(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
