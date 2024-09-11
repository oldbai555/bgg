/**
 * @Author: zjj
 * @Date: 2024/8/22
 * @Desc:
**/

package wsmgr

import "github.com/oldbai555/bgg/singlesrv/client"

func PacketWebsocketDataByChatMessage(chatMsg *client.ModelMessage) *client.WebsocketMsg {
	return &client.WebsocketMsg{
		Type:        uint32(client.WebsocketDataType_WebsocketDataTypeChat),
		ChatMessage: chatMsg,
	}
}

func PacketWebsocketDataByLogout() *client.WebsocketMsg {
	return &client.WebsocketMsg{
		Type: uint32(client.WebsocketDataType_WebsocketDataTypeLogout),
	}
}

func PacketWebsocketDataByJoinChatRoom(data *client.JoinChatRoom) *client.WebsocketMsg {
	return &client.WebsocketMsg{
		Type:         uint32(client.WebsocketDataType_WebsocketDataTypeJoinChatRoom),
		JoinChatRoom: data,
	}
}

func PacketWebsocketDataByLeaveChatRoom(data *client.LeaveChatRoom) *client.WebsocketMsg {
	return &client.WebsocketMsg{
		Type:          uint32(client.WebsocketDataType_WebsocketDataTypeLeaveChatRoom),
		LeaveChatRoom: data,
	}
}
