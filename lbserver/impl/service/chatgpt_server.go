package service

import (
	"context"
	"github.com/oldbai555/bgg/client/lbchatgpt"
)

var ChatgptServer LbchatgptServer

type LbchatgptServer struct {
	*lbchatgpt.UnimplementedLbchatgptServer
}

func (a *LbchatgptServer) ChatCompletion(ctx context.Context, req *lbchatgpt.ChatCompletionReq) (*lbchatgpt.ChatCompletionRsp, error) {
	var rsp lbchatgpt.ChatCompletionRsp
	var _ error

	return &rsp, nil
}
