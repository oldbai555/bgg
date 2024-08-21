package server

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/lbtool/pkg/routine"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/oldbai555/lbtool/log"
)

type wsConn struct {
	ws         *websocket.Conn
	writerChan chan []byte
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理WebSocket升级和管理连接。
func handleWs(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("up grade err:%v", err)
		return
	}
	c := &wsConn{
		writerChan: make(chan []byte, 512),
		ws:         ws,
	}
	routine.GoV2(func() error {
		c.writer()
		return nil
	})
	c.reader()
}

// writer 向客户端发送消息。
func (c *wsConn) writer() {
	for data := range c.writerChan {
		if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Errorf("writer err:%v", err)
			return
		}
	}
	err := c.ws.Close()
	if err != nil {
		log.Errorf("close err: %v", err)
	}
}

// reader 从客户端读取消息。
func (c *wsConn) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Errorf("read message err:%v", err)
			break
		}
		var data client.WebsocketMsg
		if err := jsonpb.Unmarshal(message, &data); err != nil {
			log.Errorf("unmarshal err:%v", err)
			continue
		}
		handleMessage(c, &data)
	}
}

// handleMessageType 处理不同类型的传入消息。
func handleMessage(c *wsConn, data *client.WebsocketMsg) {
	switch data.Type {
	case uint32(client.WebsocketDataType_WebsocketDataTypeLogin):

	case uint32(client.WebsocketDataType_WebsocketDataTypeLogout):

	case uint32(client.WebsocketDataType_WebsocketDataTypeChat):
		bytes, err := jsonpb.Marshal(data)
		if err != nil {
			log.Errorf("unmarshal err:%v", err)
			return
		}
		c.writerChan <- bytes
	default:
		log.Errorf("unknown message type:%d", data.Type)
	}
}
