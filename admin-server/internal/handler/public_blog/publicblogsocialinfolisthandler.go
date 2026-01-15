// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package public_blog

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/public_blog"
	"postapocgame/admin-server/internal/svc"
)

func PublicBlogSocialInfoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public_blog.NewPublicBlogSocialInfoListLogic(r.Context(), svcCtx)
		resp, err := l.PublicBlogSocialInfoList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
