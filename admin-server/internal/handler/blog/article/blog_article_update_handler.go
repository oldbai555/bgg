// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/blog/article"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func BlogArticleUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlogArticleUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := article.NewBlogArticleUpdateLogic(r.Context(), svcCtx)
		resp, err := l.BlogArticleUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
