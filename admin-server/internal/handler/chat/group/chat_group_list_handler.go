// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/chat/group"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func ChatGroupListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatGroupListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewChatGroupListLogic(r.Context(), svcCtx)
		resp, err := l.ChatGroupList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
