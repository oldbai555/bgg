/**
 * @Author: zjj
 * @Date: 2024/8/21
 * @Desc:
**/

package server

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
)

func handleWebsocketDataTypeLogin(ctx uctx.IUCtx, msg *client.WebsocketMsg) (proto.Message, error) {
	ctx.SetSid(msg.Sid)
	_, err := CheckAuth(ctx)
	if err != nil {
		return nil, err
	}
	wsConnMgr.bindUser2Conn(msg.Sid, ctx.ExtInfo().(string))
	return msg, nil
}

func handleWebsocketDataTypeLogout(ctx uctx.IUCtx, msg *client.WebsocketMsg) (proto.Message, error) {
	wsConnMgr.delConn(ctx.ExtInfo().(string))
	return msg, nil
}

func handleWebsocketDataTypeChat(ctx uctx.IUCtx, msg *client.WebsocketMsg) (proto.Message, error) {
	return msg, nil
}

func init() {
	regWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeLogin), handleWebsocketDataTypeLogin)
	regWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeLogout), handleWebsocketDataTypeLogout)
	regWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeChat), handleWebsocketDataTypeChat)
}
