package chat

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/pkg/response"
	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/chatclient"

	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源（生产环境应该限制）
		return true
	},
}

// ChatWSHandler 是 16-rpc-conventions.md 第 7 节 WS<->gRPC 双向流桥接骨架的网关侧实现：
// gateway 继续终结 WebSocket 连接（唯一有公网端口的进程不变），鉴权/黑名单检查（依赖 IAM
// 的 token blacklist，物理上还在这个进程里）继续留在这里，鉴权通过后建立一条到 chat-rpc
// 的 gRPC 双向流，两个方向各一个 goroutine 转发帧，chat-rpc 侧的 ChatHub 连接表和实际推送
// 逻辑见 services/chat/internal/hub、services/chat/internal/logic/streamlogic.go。
func ChatWSHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从查询参数获取 token（WebSocket 无法使用 Authorization header）
		token := r.URL.Query().Get("token")
		if token == "" {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					token = parts[1]
				}
			}
		}

		if token == "" {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "未提供认证信息"))
			return
		}

		claims, err := jwthelper.ParseToken(token, svcCtx.Config.JWT.AccessSecret)
		if err != nil || claims.IsRefresh {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "访问令牌无效或已过期"))
			return
		}

		// 检查黑名单：直连共享 Redis，key 格式与 services/iam 内部的
		// TokenBlacklistRepository 保持一致，见 16-rpc-conventions.md 第 6 节"直接复制不共享"。
		blacklisted, err := svcCtx.Redis.Exists(consts.RedisJWTBlacklistPrefix + token)
		if err != nil {
			response.ErrorCtx(r.Context(), w, errs.Wrap(errs.CodeInternalError, "检查令牌黑名单失败", err))
			return
		}
		if blacklisted {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "令牌已失效"))
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Errorf("WebSocket 升级失败: %v", err)
			return
		}
		defer conn.Close()

		stream, err := svcCtx.ChatRPC.Stream(r.Context())
		if err != nil {
			logx.Errorf("建立 chat-rpc 流失败: userId=%d, err=%v", claims.UserID, err)
			return
		}

		if err := stream.Send(&chatclient.ClientFrame{Payload: &chat.ClientFrame_Join{
			Join: &chatclient.JoinFrame{UserId: claims.UserID, Username: claims.Username},
		}}); err != nil {
			logx.Errorf("发送 JoinFrame 失败: userId=%d, err=%v", claims.UserID, err)
			return
		}

		errCh := make(chan error, 2)

		// WS -> gRPC：读浏览器发来的帧，转成 ClientFrame 发给 chat-rpc。当前真实前端只用
		// WS 做服务端推送、发消息走 REST（见 services/chat/rpc/chat.proto 顶部注释），这里
		// 按文档骨架实现完整但未必有真实调用方在用的路径。
		go func() {
			for {
				_, data, err := conn.ReadMessage()
				if err != nil {
					errCh <- err
					return
				}
				frame, ok := decodeClientFrame(data)
				if !ok {
					continue // 单帧解析失败不断连，忽略即可
				}
				if err := stream.Send(frame); err != nil {
					errCh <- err
					return
				}
			}
		}()

		// gRPC -> WS：读 chat-rpc 推来的 ServerFrame，原样把 MessageFrame.PayloadJson 的
		// 字节写回浏览器（保持拆分前的 WS wire 格式不变，见 chat.proto MessageFrame 注释）。
		go func() {
			for {
				frame, err := stream.Recv()
				if err != nil {
					errCh <- err
					return
				}
				data, ok := encodeServerFrame(frame)
				if !ok {
					continue
				}
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					errCh <- err
					return
				}
			}
		}()

		<-errCh
	}
}

// decodeClientFrame 把浏览器发来的 JSON 帧转成 pb.ClientFrame。当前唯一有意义的入站帧是
// 心跳（{"type":"ping"}），发消息走 REST 不走这条路径（见上方注释）。
func decodeClientFrame(data []byte) (*chatclient.ClientFrame, bool) {
	var raw struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, false
	}
	if raw.Type == "ping" {
		return &chatclient.ClientFrame{Payload: &chat.ClientFrame_Ping{Ping: &chatclient.PingFrame{}}}, true
	}
	return nil, false
}

// encodeServerFrame 把 chat-rpc 推来的 ServerFrame 转成要写回 WS 连接的字节。MessageFrame
// 原样透传 PayloadJson（保持 hub.ChatMessage 的 WS wire 格式不变）；Pong/Error 两种协议级
// 帧编码成最小 JSON 信封，当前没有真实前端消费这两种，仅为文档骨架完整性保留。
func encodeServerFrame(frame *chatclient.ServerFrame) ([]byte, bool) {
	switch p := frame.Payload.(type) {
	case *chat.ServerFrame_Message:
		return []byte(p.Message.PayloadJson), true
	case *chat.ServerFrame_Pong:
		return []byte(`{"type":"pong"}`), true
	case *chat.ServerFrame_Error:
		b, _ := json.Marshal(map[string]string{"type": "error", "message": p.Error.Message})
		return b, true
	default:
		return nil, false
	}
}
