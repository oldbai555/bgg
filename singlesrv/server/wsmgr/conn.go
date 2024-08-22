/**
 * @Author: zjj
 * @Date: 2024/8/22
 * @Desc:
**/

package wsmgr

import (
	"github.com/gorilla/websocket"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"sync/atomic"
)

type wsConn struct {
	ws         *websocket.Conn
	writerChan chan []byte
	connId     string
	uid        uint64
	close      atomic.Bool
}

func (c *wsConn) GetConnId() string {
	return c.connId
}

func (c *wsConn) GetUid() uint64 {
	return c.uid
}

func (c *wsConn) IsLogin() bool {
	return c.uid != 0
}

// writer 向客户端发送消息。
func (c *wsConn) writer() {
	for data := range c.writerChan {
		if c.close.Load() {
			break
		}
		if err := c.ws.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Errorf("writer err:%v", err)
			break
		}
	}
}

// reader 从客户端读取消息。
func (c *wsConn) reader() {
	for {
		if c.close.Load() {
			break
		}
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Errorf("read err:%v", err)
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

func (c *wsConn) writeDisConnectMsg() {
	val := &client.WebsocketMsg{Type: uint32(client.WebsocketDataType_WebsocketDataTypeDisConnect)}
	bytes, err := jsonpb.Marshal(val)
	if err != nil {
		log.Errorf("unmarshal err:%v", err)
		return
	}
	err = c.ws.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Errorf("write err:%v", err)
	}
}
