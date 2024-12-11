/**
 * @Author: zjj
 * @Date: 2024/8/22
 * @Desc:
**/

package wsmgr

import (
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingle"
)

func PacketWebsocketDataByChatMessage(chatMsg *lbsingle.ModelMessage) *lbsingle.WebsocketMsg {
	return &lbsingle.WebsocketMsg{
		Type:        uint32(lbbase.WebsocketDataType_WebsocketDataTypeChat),
		ChatMessage: chatMsg,
	}
}

func PacketWebsocketDataByLogout() *lbsingle.WebsocketMsg {
	return &lbsingle.WebsocketMsg{
		Type: uint32(lbbase.WebsocketDataType_WebsocketDataTypeLogout),
	}
}

func PacketWebsocketDataByJoinChatRoom(data *lbbase.JoinChatRoom) *lbsingle.WebsocketMsg {
	return &lbsingle.WebsocketMsg{
		Type:         uint32(lbbase.WebsocketDataType_WebsocketDataTypeJoinChatRoom),
		JoinChatRoom: data,
	}
}

func PacketWebsocketDataByLeaveChatRoom(data *lbbase.LeaveChatRoom) *lbsingle.WebsocketMsg {
	return &lbsingle.WebsocketMsg{
		Type:          uint32(lbbase.WebsocketDataType_WebsocketDataTypeLeaveChatRoom),
		LeaveChatRoom: data,
	}
}
