// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article_audit

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/blog/article_audit"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func BlogArticleAuditUnpublishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlogArticleAuditUnpublishReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := article_audit.NewBlogArticleAuditUnpublishLogic(r.Context(), svcCtx)
		resp, err := l.BlogArticleAuditUnpublish(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
