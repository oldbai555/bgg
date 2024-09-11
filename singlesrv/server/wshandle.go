/**
 * @Author: zjj
 * @Date: 2024/8/21
 * @Desc:
**/

package server

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/iface"
	"github.com/oldbai555/bgg/singlesrv/server/roommgr"
	"github.com/oldbai555/bgg/singlesrv/server/wsmgr"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
)

func handleWebsocketDataTypeLogin(ctx uctx.IUCtx, msg *client.WebsocketMsg) (ret proto.Message, err error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	if conn.IsLogin() {
		return
	}

	login := msg.Login
	ctx.SetSid(login.Sid)

	// 鉴权
	info, err := CheckAuth(ctx)
	if err != nil {
		return nil, err
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
	newMsg := wsmgr.PacketWebsocketDataByJoinChatRoom(&client.JoinChatRoom{RoomId: 1, Member: &client.ChatRoomMember{Uid: uid, Username: info.Username}})
	err = chatRoomSt.Broadcast(newMsg)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	return
}

func handleWebsocketDataTypeLogout(ctx uctx.IUCtx, msg *client.WebsocketMsg) (ret proto.Message, err error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	if !conn.IsLogin() {
		return
	}
	uid := conn.GetUid()

	newMsg := wsmgr.PacketWebsocketDataByLeaveChatRoom(&client.LeaveChatRoom{RoomId: 1, Member: &client.ChatRoomMember{Uid: uid}})
	err = roommgr.GetSingleChatRoom().Broadcast(newMsg)
	if err != nil {
		return
	}

	// 绑定用户和链接
	wsmgr.UnBindUser2Conn(uid)
	wsmgr.WriteProtoMsg(uid, msg)
	roommgr.GetSingleChatRoom().DelUser(uid)
	return
}

func handleWebsocketDataTypeChat(ctx uctx.IUCtx, msg *client.WebsocketMsg) (ret proto.Message, err error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	if !conn.IsLogin() {
		return
	}

	err = roommgr.GetSingleChatRoom().Broadcast(msg)
	if err != nil {
		return
	}

	return
}

func handleWebsocketDataTypeHeartBeat(_ uctx.IUCtx, msg *client.WebsocketMsg) (ret proto.Message, err error) {
	ret = msg
	return
}

func handleWebsocketDataTypeConnect(_ uctx.IUCtx, msg *client.WebsocketMsg) (ret proto.Message, err error) {
	ret = msg
	return
}

func handleWebsocketDataTypeDisConnect(ctx uctx.IUCtx, msg *client.WebsocketMsg) (ret proto.Message, err error) {
	conn := ctx.ExtInfo().(iface.IWsConn)
	wsmgr.DelConn(conn.GetConnId())
	ret = msg
	return
}

func init() {
	wsmgr.RegWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeConnect), handleWebsocketDataTypeConnect)
	wsmgr.RegWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeDisConnect), handleWebsocketDataTypeDisConnect)
	wsmgr.RegWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeHeartBeat), handleWebsocketDataTypeHeartBeat)
	wsmgr.RegWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeLogin), handleWebsocketDataTypeLogin)
	wsmgr.RegWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeLogout), handleWebsocketDataTypeLogout)
	wsmgr.RegWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeChat), handleWebsocketDataTypeChat)
}
