package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/ctx"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/oldbai555/lbtool/log"
)

type wsConn struct {
	ws         *websocket.Conn
	writerChan chan []byte
	connId     string
	sid        string
	close      atomic.Bool
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
		connId:     utils.GenUUID(),
	}
	wsConnMgr.connMgr[c.connId] = c
	routine.GoV2(func() error {
		c.writer()
		return nil
	})
	c.reader()
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

// handleMessageType 处理不同类型的传入消息。
func handleMessage(c *wsConn, data *client.WebsocketMsg) {
	f, ok := wsMsgTypeHandleMgr[data.Type]
	if !ok {
		log.Errorf("unknown message type:%d", data.Type)
		return
	}
	nCtx := ctx.NewCtx(context.Background())
	nCtx.SetExtInfo(c.connId)
	val, err := f(nCtx, data)
	if val == nil && err == nil {
		return
	}
	if err != nil {
		val = &client.WebsocketMsg{ErrMsg: err.Error()}
	}
	bytes, err := jsonpb.Marshal(val)
	if err != nil {
		log.Errorf("unmarshal err:%v", err)
		return
	}
	c.writerChan <- bytes
}

type wsMsgTypeHandleFunc func(ctx uctx.IUCtx, msg *client.WebsocketMsg) (proto.Message, error)

var wsMsgTypeHandleMgr = make(map[uint32]wsMsgTypeHandleFunc)

func regWsMsgTypeHandler(msgType uint32, handler wsMsgTypeHandleFunc) {
	_, ok := wsMsgTypeHandleMgr[msgType]
	if ok {
		panic(fmt.Sprintf("msgTypeHandler already exist %d", msgType))
		return
	}
	wsMsgTypeHandleMgr[msgType] = handler
}

type wsConnMgrSt struct {
	connMgr   map[string]*wsConn // key by connId
	user2Conn map[string]*wsConn // key by sid
}

func (s *wsConnMgrSt) bindUser2Conn(sid string, connId string) {
	if len(sid) == 0 || len(connId) == 0 {
		return
	}
	conn, ok := s.connMgr[connId]
	if !ok {
		return
	}
	conn.sid = sid
	s.user2Conn[sid] = conn
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

func (s *wsConnMgrSt) delConn(connId string) {
	if len(connId) == 0 {
		return
	}
	conn, ok := s.user2Conn[connId]
	if !ok {
		return
	}
	delete(s.user2Conn, conn.sid)
	delete(s.connMgr, conn.connId)
	conn.close.Store(true)
}

var wsConnMgr = &wsConnMgrSt{
	connMgr:   make(map[string]*wsConn),
	user2Conn: make(map[string]*wsConn),
}

func CloseAllWsConn() {
	for _, conn := range wsConnMgr.connMgr {
		wsConnMgr.delConn(conn.connId)
		conn.writeDisConnectMsg()
	}
}

func init() {
	regWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeConnect), handleWebsocketDataTypeConnect)
	regWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeDisConnect), handleWebsocketDataTypeDisConnect)
	regWsMsgTypeHandler(uint32(client.WebsocketDataType_WebsocketDataTypeHeartBeat), handleWebsocketDataTypeHeartBeat)
}

func handleWebsocketDataTypeHeartBeat(_ uctx.IUCtx, msg *client.WebsocketMsg) (proto.Message, error) {
	return msg, nil
}

func handleWebsocketDataTypeConnect(_ uctx.IUCtx, msg *client.WebsocketMsg) (proto.Message, error) {
	return msg, nil
}

func handleWebsocketDataTypeDisConnect(ctx uctx.IUCtx, msg *client.WebsocketMsg) (proto.Message, error) {
	wsConnMgr.delConn(ctx.ExtInfo().(string))
	return msg, nil
}
