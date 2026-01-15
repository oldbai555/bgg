// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_social_info

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/blog_social_info"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func BlogSocialInfoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlogSocialInfoListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := blog_social_info.NewBlogSocialInfoListLogic(r.Context(), svcCtx)
		resp, err := l.BlogSocialInfoList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
