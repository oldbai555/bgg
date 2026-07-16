package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type PushToUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushToUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushToUserLogic {
	return &PushToUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PushToUser 供还留在单体里的 internal/consumer/task_notification_consumer.go 回调，把
// 已经拼好的 hub.ChatMessage JSON 原样转发给目标用户的在线连接。见 chat.proto 里
// PushToUserRequest 的注释。
func (l *PushToUserLogic) PushToUser(in *chat.PushToUserRequest) (*chat.PushToUserResponse, error) {
	frame := &chat.ServerFrame{Payload: &chat.ServerFrame_Message{
		Message: &chat.MessageFrame{PayloadJson: in.PayloadJson},
	}}
	delivered := l.svcCtx.Hub.SendToUser(in.UserId, frame)
	if !delivered {
		logx.Infof("PushToUser 用户不在线，消息丢弃: userId=%d", in.UserId)
	}
	return &chat.PushToUserResponse{Delivered: delivered}, nil
}
