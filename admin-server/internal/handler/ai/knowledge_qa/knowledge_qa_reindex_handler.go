// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package knowledge_qa

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/ai/knowledge_qa"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
)

func KnowledgeQaReindexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.KnowledgeQaReindexReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := knowledge_qa.NewKnowledgeQaReindexLogic(r.Context(), svcCtx)
		resp, err := l.KnowledgeQaReindex(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
