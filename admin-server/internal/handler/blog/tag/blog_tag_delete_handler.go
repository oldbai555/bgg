// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tag

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/blog/tag"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func BlogTagDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlogTagDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tag.NewBlogTagDeleteLogic(r.Context(), svcCtx)
		err := l.BlogTagDelete(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
