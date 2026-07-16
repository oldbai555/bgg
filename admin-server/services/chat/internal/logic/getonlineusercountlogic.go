package logic

import (
	"context"

	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineUserCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOnlineUserCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineUserCountLogic {
	return &GetOnlineUserCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOnlineUserCount 取代原来 gateway 直接读 svcCtx.ChatHub.GetOnlineUsers() 再取 len 的写法
// （internal/logic/monitoring/{monitor,login_log}/*_stats_logic.go 两处调用点）。
func (l *GetOnlineUserCountLogic) GetOnlineUserCount(in *chat.Empty) (*chat.GetOnlineUserCountResponse, error) {
	return &chat.GetOnlineUserCountResponse{Count: l.svcCtx.Hub.OnlineUserCount()}, nil
}
