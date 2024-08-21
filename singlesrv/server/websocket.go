package server

import (
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/lbtool/pkg/routine"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oldbai555/lbtool/log"
)

type wsConn struct {
	ws       *websocket.Conn
	dataChan chan []byte
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWs 处理WebSocket升级和管理连接。
func HandleWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("WebSocket升级错误: %v", err)
		return
	}
	c := &wsConn{
		dataChan: make(chan []byte, 512),
		ws:       ws,
	}
	routine.GoV2(func() error {
		c.writer()
		return nil
	})
	c.reader()
}

// writer 向客户端发送消息。
func (c *wsConn) writer() {
	for data := range c.dataChan {
		if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Errorf("写入错误: %v", err)
			return
		}
	}
	err := c.ws.Close()
	if err != nil {
		log.Errorf("关闭错误: %v", err)
	}
}

// reader 从客户端读取消息。
func (c *wsConn) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Errorf("读取错误: %v", err)
			break
		}
		var data client.WebsocketData
		if err := jsonpb.Unmarshal(message, &data); err != nil {
			log.Errorf("反序列化错误: %v", err)
			continue
		}
		handleMessage(c, &data)
	}
}

// handleMessageType 处理不同类型的传入消息。
func handleMessage(c *wsConn, data *client.WebsocketData) {
	switch data.Type {
	case uint32(client.WebsocketDataType_WebsocketDataTypeLogin):

	case uint32(client.WebsocketDataType_WebsocketDataTypeLogout):

	case uint32(client.WebsocketDataType_WebsocketDataTypeChat):

	default:
		log.Errorf("无法识别消息类型%d", data.Type)
	}
}
