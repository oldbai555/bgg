// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package task_public

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/task_public"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func TaskRecentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TaskRecentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := task_public.NewTaskRecentLogic(r.Context(), svcCtx)
		resp, err := l.TaskRecent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
