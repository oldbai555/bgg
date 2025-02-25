/**
 * @Author: zjj
 * @Date: 2024/8/21
 * @Desc:
**/

package lbsingleserver

import (
	"github.com/oldbai555/bgg/iface"
	"github.com/oldbai555/bgg/pkg/syscfg"
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/deepseek"
	"github.com/oldbai555/bgg/service/lbsingleserver/roommgr"
	"github.com/oldbai555/bgg/service/lbsingleserver/wsmgr"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
	"time"
)

func handleWebsocketDataTypeLogin(ctx uctx.IUCtx, msg *lbsingle.WebsocketMsg) (ret proto.Message, err error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	if conn.IsLogin() {
		return
	}

	login := msg.Login
	ctx.SetSid(login.Sid)

	// 鉴权
	info, err := CheckAuth(ctx)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	uid := info.Id

	// 绑定用户和链接
	wsmgr.BindUser2Conn(uid, conn.GetConnId())

	// 加入全局房间
	chatRoomSt := roommgr.GetSingleChatRoom()
	chatRoomSt.AddUser(uid)

	login.Uid = uid

	// 回复登陆成功
	wsmgr.WriteProtoMsg(uid, msg)

	// 广播有人加入房间
	newMsg := wsmgr.PacketWebsocketDataByJoinChatRoom(&lbbase.JoinChatRoom{RoomId: 1, Member: &lbbase.ChatRoomMember{Uid: uid, Username: info.Username}})
	err = chatRoomSt.Broadcast(newMsg)
	if err != nil {
		err = lberr.Wrap(err)
		return
	}
	return
}

func handleWebsocketDataTypeLogout(ctx uctx.IUCtx, msg *lbsingle.WebsocketMsg) (ret proto.Message, err error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	if !conn.IsLogin() {
		return
	}
	uid := conn.GetUid()

	newMsg := wsmgr.PacketWebsocketDataByLeaveChatRoom(&lbbase.LeaveChatRoom{RoomId: 1, Member: &lbbase.ChatRoomMember{Uid: uid}})
	err = roommgr.GetSingleChatRoom().Broadcast(newMsg)
	if err != nil {
		err = lberr.Wrap(err)
		return
	}

	// 绑定用户和链接
	wsmgr.UnBindUser2Conn(uid)
	wsmgr.WriteProtoMsg(uid, msg)
	roommgr.GetSingleChatRoom().DelUser(uid)
	return
}

var cacheMsgList []*deepseek.ChatCompletionsMessage
var defaultFirstMsg = &deepseek.ChatCompletionsMessage{
	Content: "请你扮演一个无所不知的学者,说话通俗易懂,面对我们提出的问题你都能一一回答解决",
	Role:    "system",
}

func appendCacheMsgList(msg *deepseek.ChatCompletionsMessage) {
	cacheMsgList = append(cacheMsgList, msg)
	if len(cacheMsgList) >= 30 {
		cacheMsgList = cacheMsgList[1:]
	}
}

func buildMsgList(msg *deepseek.ChatCompletionsMessage) []*deepseek.ChatCompletionsMessage {
	var msgList []*deepseek.ChatCompletionsMessage
	msgList = append(msgList, defaultFirstMsg)
	msgList = append(msgList, cacheMsgList...)
	msgList = append(msgList, msg)
	appendCacheMsgList(msg)
	return msgList
}

func handleWebsocketDataTypeChat(ctx uctx.IUCtx, msg *lbsingle.WebsocketMsg) (proto.Message, error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	if !conn.IsLogin() {
		return nil, lbsingle.ErrUserNotFound
	}

	// 先广播给房间的所有人
	err := roommgr.GetSingleChatRoom().Broadcast(msg)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	// 再通知deepseek去生成回答
	if msg.ChatMessage.To == "deepseek" {
		deepSeekCfg, err := syscfg.GetDeepSeek()
		if err != nil {
			return nil, lberr.Wrap(err)
		}
		completions, err := deepseek.ChatCompletions(&deepseek.ChatChatCompletionsReq{
			BaseUrl: deepSeekCfg.BaseUrl,
			Token:   deepSeekCfg.Token,
			MsgList: buildMsgList(&deepseek.ChatCompletionsMessage{
				Content: msg.ChatMessage.Content.Text.Content,
				Role:    "user",
			}),
		})
		for _, message := range completions.MsgList {
			appendCacheMsgList(message)
			err = roommgr.GetSingleChatRoom().Broadcast(wsmgr.PacketWebsocketDataByChatMessage(&lbsingle.ModelMessage{
				SendAt: uint64(time.Now().UnixMilli()),
				From:   "deepseek",
				To:     msg.ChatMessage.From,
				Content: &lbbase.Content{
					Text: &lbbase.Content_Text{
						Content: message.Content,
					},
				},
			}))
			if err != nil {
				return nil, lberr.Wrap(err)
			}
		}
	}

	return nil, nil
}

func handleWebsocketDataTypeHeartBeat(_ uctx.IUCtx, msg *lbsingle.WebsocketMsg) (ret proto.Message, err error) {
	ret = msg
	return
}

func handleWebsocketDataTypeConnect(_ uctx.IUCtx, msg *lbsingle.WebsocketMsg) (ret proto.Message, err error) {
	ret = msg
	return
}

func handleWebsocketDataTypeDisConnect(ctx uctx.IUCtx, msg *lbsingle.WebsocketMsg) (ret proto.Message, err error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	wsmgr.DelConn(conn.GetConnId())
	ret = msg
	return
}

func init() {
	wsmgr.RegWsMsgTypeHandler(uint32(lbbase.WebsocketDataType_WebsocketDataTypeConnect), handleWebsocketDataTypeConnect)
	wsmgr.RegWsMsgTypeHandler(uint32(lbbase.WebsocketDataType_WebsocketDataTypeDisConnect), handleWebsocketDataTypeDisConnect)
	wsmgr.RegWsMsgTypeHandler(uint32(lbbase.WebsocketDataType_WebsocketDataTypeHeartBeat), handleWebsocketDataTypeHeartBeat)
	wsmgr.RegWsMsgTypeHandler(uint32(lbbase.WebsocketDataType_WebsocketDataTypeLogin), handleWebsocketDataTypeLogin)
	wsmgr.RegWsMsgTypeHandler(uint32(lbbase.WebsocketDataType_WebsocketDataTypeLogout), handleWebsocketDataTypeLogout)
	wsmgr.RegWsMsgTypeHandler(uint32(lbbase.WebsocketDataType_WebsocketDataTypeChat), handleWebsocketDataTypeChat)
}
