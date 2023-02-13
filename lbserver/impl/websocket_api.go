package impl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/lbwebsocket"
	"github.com/oldbai555/lbtool/log"
)

func registerLbwebsocketApi(h *gin.Engine) {
	// 可以利用反射来映射函数进去
	group := h.Group("ws")
	group.GET("/:vid", HandleWs)
}

func HandleWs(c *gin.Context) {
	var req lbwebsocket.HandleWsReq
	handler := NewHandler(c)

	wsConn, err := wu.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("err is %v", err)
		handler.RespError(err)
		return
	}

	conn, err := InitConnection(wsConn)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	data, err := conn.ReadMessage()
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	var wsData lbwebsocket.WsData
	err = json.Unmarshal(data, &wsData)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	wsData.Vid = handler.C.Param("vid")
	req.WsData = &wsData
	req.Ip = c.ClientIP()

	// 数据入库
	rsp, err := lbwebsocketServer.HandleWs(c, &req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	marshal, err := json.Marshal(rsp.WsData)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	err = conn.WriteMessage(marshal)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}
}
