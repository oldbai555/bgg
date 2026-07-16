// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package monitor

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"postapocgame/admin-server/internal/logic/monitoring/monitor"
	"postapocgame/admin-server/internal/svc"
)

func MonitorStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := monitor.NewMonitorStatsLogic(r.Context(), svcCtx)
		resp, err := l.MonitorStats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
