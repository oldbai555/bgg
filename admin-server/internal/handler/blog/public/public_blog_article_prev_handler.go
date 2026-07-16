// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/blog/public"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func PublicBlogArticlePrevHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublicBlogArticlePrevReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := public.NewPublicBlogArticlePrevLogic(r.Context(), svcCtx)
		resp, err := l.PublicBlogArticlePrev(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
