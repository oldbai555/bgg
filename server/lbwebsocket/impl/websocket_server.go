package impl

import (
	"context"
	"github.com/oldbai555/bgg/client/lbwebsocket"
)

var lbwebsocketServer LbwebsocketServer

type LbwebsocketServer struct {
	*lbwebsocket.UnimplementedLbwebsocketServer
}

func (a *LbwebsocketServer) HandleWs(ctx context.Context, req *lbwebsocket.HandleWsReq) (*lbwebsocket.HandleWsRsp, error) {
	var rsp lbwebsocket.HandleWsRsp
	return &rsp, nil
}
