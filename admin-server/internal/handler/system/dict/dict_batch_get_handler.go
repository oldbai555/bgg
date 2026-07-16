// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/system/dict"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func DictBatchGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DictBatchGetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dict.NewDictBatchGetLogic(r.Context(), svcCtx)
		resp, err := l.DictBatchGet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
