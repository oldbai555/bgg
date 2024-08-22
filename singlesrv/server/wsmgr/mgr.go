package wsmgr

import (
	"github.com/gorilla/websocket"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/iface"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsConnMgrSt struct {
	connMgr   map[string]*wsConn // key by connId
	user2Conn map[uint64]*wsConn // key by sid
	rwLock    sync.RWMutex
}

func (s *wsConnMgrSt) bindUser2Conn(uid uint64, connId string) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	if uid == 0 || len(connId) == 0 {
		return
	}
	conn, ok := s.connMgr[connId]
	if !ok {
		return
	}
	conn.uid = uid
	s.user2Conn[uid] = conn
}
func (s *wsConnMgrSt) unBindUser2Conn(uid uint64) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	conn, ok := s.user2Conn[uid]
	if !ok {
		return
	}
	conn.uid = 0
	delete(s.user2Conn, uid)
}

func (s *wsConnMgrSt) delConn(connId string) {
	wsConnMgr.rwLock.Lock()
	defer wsConnMgr.rwLock.Unlock()
	if len(connId) == 0 {
		return
	}
	conn, ok := s.connMgr[connId]
	if !ok {
		return
	}
	delete(s.user2Conn, conn.uid)
	delete(s.connMgr, conn.connId)
	conn.uid = 0
	conn.close.Store(true)
}

func (s *wsConnMgrSt) writeError(uid uint64, err error) {
	val := &client.WebsocketMsg{Type: uint32(client.WebsocketDataType_WebsocketDataTypeError), ErrMsg: err.Error()}
	s.writeProtoMsg(uid, val)
}

func (s *wsConnMgrSt) writeProtoMsg(uid uint64, msg proto.Message) {
	bytes, err := jsonpb.Marshal(msg)
	if err != nil {
		log.Errorf("unmarshal err:%v", err)
		return
	}
	s.writeBytes(uid, bytes)
}

func (s *wsConnMgrSt) writeBytes(uid uint64, bytes []byte) {
	wsConnMgr.rwLock.RLock()
	defer wsConnMgr.rwLock.RUnlock()
	c, ok := s.user2Conn[uid]
	if !ok {
		return
	}
	err := c.ws.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Errorf("write err:%v", err)
		return
	}
}

func (s *wsConnMgrSt) closeConnByUid(uid uint64) {
	s.rwLock.RLock()
	conn, ok := s.user2Conn[uid]
	s.rwLock.RUnlock()
	if !ok {
		return
	}
	conn.writeDisConnectMsg()
	conn.close.Store(true)
	wsConnMgr.delConn(conn.connId)
}

func (s *wsConnMgrSt) getConnByUid(uid uint64) iface.IWsConn {
	s.rwLock.RLock()
	conn, ok := s.user2Conn[uid]
	s.rwLock.RUnlock()
	if !ok {
		return nil
	}
	return conn
}

var wsConnMgr = &wsConnMgrSt{
	connMgr:   make(map[string]*wsConn),
	user2Conn: make(map[uint64]*wsConn),
}
