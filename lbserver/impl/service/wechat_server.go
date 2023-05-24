package service

import (
	"context"
	"github.com/oldbai555/bgg/client/lbwechat"
)

var WechatServer LbwechatServer

type LbwechatServer struct {
	*lbwechat.UnimplementedLbwechatServer
}

func (a *LbwechatServer) HandleWxGzhAuth(ctx context.Context, req *lbwechat.HandleWxGzhAuthReq) (*lbwechat.HandleWxGzhAuthRsp, error) {
	var rsp lbwechat.HandleWxGzhAuthRsp
	var _ error

	return &rsp, nil
}
func (a *LbwechatServer) HandleWxGzhMsg(ctx context.Context, req *lbwechat.HandleWxGzhMsgReq) (*lbwechat.HandleWxGzhMsgRsp, error) {
	var rsp lbwechat.HandleWxGzhMsgRsp
	var _ error

	return &rsp, nil
}
