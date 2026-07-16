// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict_item

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/system/dict_item"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func DictItemCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DictItemCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dict_item.NewDictItemCreateLogic(r.Context(), svcCtx)
		err := l.DictItemCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
