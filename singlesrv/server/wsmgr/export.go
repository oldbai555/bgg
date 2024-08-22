/**
 * @Author: zjj
 * @Date: 2024/8/22
 * @Desc:
**/

package wsmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
)

var (
	DelConn         = wsConnMgr.delConn
	BindUser2Conn   = wsConnMgr.bindUser2Conn
	UnBindUser2Conn = wsConnMgr.unBindUser2Conn
	WriteProtoMsg   = wsConnMgr.writeProtoMsg
	WriteBytes      = wsConnMgr.writeBytes
)

// HandleWs 处理WebSocket升级和管理连接。
func HandleWs(ctx *gin.Context) {
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

	// 读取失败后就开始关闭连接
	wsConnMgr.delConn(c.connId)
	err = c.ws.Close()
	if err != nil {
		log.Errorf("up grade err:%v", err)
	}
}

func CloseAllWsConn() {
	var connIdsToDelete []string
	wsConnMgr.rwLock.RLock()
	for _, conn := range wsConnMgr.connMgr {
		conn.writeDisConnectMsg()
		conn.close.Store(true)
		connIdsToDelete = append(connIdsToDelete, conn.connId)
	}
	wsConnMgr.rwLock.RUnlock()
	for _, connId := range connIdsToDelete {
		wsConnMgr.delConn(connId)
	}
}
