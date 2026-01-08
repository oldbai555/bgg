// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package m3u8

import (
	"net/http"

	"postapocgame/admin-server/internal/logic/m3u8"
	"postapocgame/admin-server/internal/svc"
)

func M3u8ProxyOptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := m3u8.NewM3u8ProxyOptionsLogic(r.Context(), svcCtx)
		l.M3u8ProxyOptions(w, r)
		// OPTIONS请求直接返回200 OK，不需要错误处理
	}
}
