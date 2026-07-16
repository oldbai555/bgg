package logic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/services/chat/chat"
	chathub "postapocgame/admin-server/services/chat/internal/hub"
	"postapocgame/admin-server/services/chat/internal/svc"
)

var errStreamFirstFrameMustBeJoin = errors.New("第一帧必须是 JoinFrame")

type StreamLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StreamLogic {
	return &StreamLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Stream 实现 16-rpc-conventions.md 第 7 节的 WS<->gRPC 双向流骨架的服务端一半：gateway
// 侧的 chatwshandler.go 每条 WebSocket 连接桥接一条到这里的 gRPC Stream，第一帧必须是
// JoinFrame（携带已经在 gateway 侧鉴权过的 user_id），之后这条连接注册进 Hub 的连接表，
// 后续的消息推送（ChatMessageSend 广播、PushToUser）通过 Hub.SendToUser/BroadcastToChat
// 塞进 client.Send，由这里的 writer goroutine 转发。
//
// 当前真实前端只用这条流做"服务端推送"，SendMessageFrame（客户端经 WS 发消息）按文档骨架
// 实现、复用与 ChatMessageSendLogic 相同的持久化+广播逻辑，但没有真实前端在用这条路径，
// 只做过 wscat 级别的手工验证，见 services/chat/rpc/chat.proto 顶部注释、
// docs/progress.md 对应条目。
func (l *StreamLogic) Stream(stream chat.Chat_StreamServer) error {
	frame, err := stream.Recv()
	if err != nil {
		return err
	}
	join := frame.GetJoin()
	if join == nil {
		return errStreamFirstFrameMustBeJoin
	}

	client := &chathub.Client{
		Hub:      l.svcCtx.Hub,
		Stream:   stream,
		Send:     make(chan *chat.ServerFrame, 256),
		UserID:   join.UserId,
		Username: join.Username,
	}
	l.svcCtx.Hub.Register() <- client
	defer func() { l.svcCtx.Hub.Unregister() <- client }()

	errCh := make(chan error, 2)

	// 写方向：Hub 推给这个连接的帧 -> stream.Send
	go func() {
		for f := range client.Send {
			if err := stream.Send(f); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// 读方向：stream.Recv -> 处理客户端帧（心跳 / 发消息）
	go func() {
		for {
			f, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			l.handleClientFrame(client, f)
		}
	}()

	return <-errCh
}

func (l *StreamLogic) handleClientFrame(client *chathub.Client, f *chat.ClientFrame) {
	switch p := f.Payload.(type) {
	case *chat.ClientFrame_Ping:
		select {
		case client.Send <- &chat.ServerFrame{Payload: &chat.ServerFrame_Pong{Pong: &chat.PongFrame{}}}:
		default:
		}
	case *chat.ClientFrame_Send:
		l.handleSendFrame(client, p.Send)
	}
}

// handleSendFrame 复用 ChatMessageSendLogic 的持久化+广播实现，operator 信息取自这条连接
// 建立时的 JoinFrame（client.UserID/Username），不是每帧都带。
func (l *StreamLogic) handleSendFrame(client *chathub.Client, f *chat.SendMessageFrame) {
	sendLogic := NewChatMessageSendLogic(l.ctx, l.svcCtx)
	if _, err := sendLogic.ChatMessageSend(&chat.ChatMessageSendRequest{
		ChatId:           f.ChatId,
		Content:          f.Content,
		MessageType:      int64(f.MessageType),
		OperatorUserId:   client.UserID,
		OperatorUsername: client.Username,
	}); err != nil {
		select {
		case client.Send <- &chat.ServerFrame{Payload: &chat.ServerFrame_Error{
			Error: &chat.ErrorFrame{Message: err.Error()},
		}}:
		default:
		}
	}
}
