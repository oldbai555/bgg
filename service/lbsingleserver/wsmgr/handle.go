/**
 * @Author: zjj
 * @Date: 2024/8/22
 * @Desc:
**/

package wsmgr

import (
	"context"
	"fmt"
	"github.com/oldbai555/bgg/iface"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/service/lbbase"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/jsonpb"
	"github.com/oldbai555/micro/uctx"
	"google.golang.org/protobuf/proto"
)

type WsMsgTypeHandleFunc func(ctx uctx.IUCtx, msg *lbsingle.WebsocketMsg) (ret proto.Message, err error)

var wsMsgTypeHandleMgr = make(map[uint32]WsMsgTypeHandleFunc)

func RegWsMsgTypeHandler(msgType uint32, handler WsMsgTypeHandleFunc) {
	_, ok := wsMsgTypeHandleMgr[msgType]
	if ok {
		panic(fmt.Sprintf("msgTypeHandler already exist %d", msgType))
		return
	}
	wsMsgTypeHandleMgr[msgType] = handler
}

// handleMessageType 处理不同类型的传入消息。
func handleMessage(c *wsConn, data *lbsingle.WebsocketMsg) {
	f, ok := wsMsgTypeHandleMgr[data.Type]
	if !ok {
		log.Errorf("unknown message type:%d", data.Type)
		return
	}

	nCtx := bctx.NewCtx(context.Background())

	// 转换成接口传下去
	var extInfo iface.IWsConn = c
	nCtx.SetExtInfo(extInfo)

	val, err := f(nCtx, data)

	// 不需要返回值
	if val == nil && err == nil {
		return
	}

	// 报错 直接替换
	if err != nil {
		val = &lbsingle.WebsocketMsg{Type: uint32(lbbase.WebsocketDataType_WebsocketDataTypeError), ErrMsg: err.Error()}
	}

	bytes, err := jsonpb.Marshal(val)
	if err != nil {
		log.Errorf("unmarshal err:%v", err)
		return
	}
	c.writerChan <- bytes
}
