// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package m3u8

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/svc"
)

type M3u8ProxyOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewM3u8ProxyOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *M3u8ProxyOptionsLogic {
	return &M3u8ProxyOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *M3u8ProxyOptionsLogic) M3u8ProxyOptions(w http.ResponseWriter, r *http.Request) {
	// 设置CORS响应头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	// 直接返回200 OK
	w.WriteHeader(http.StatusOK)
}
