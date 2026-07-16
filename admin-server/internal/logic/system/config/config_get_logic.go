// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigGetLogic {
	return &ConfigGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigGetLogic) ConfigGet(req *types.ConfigGetReq) (resp *types.ConfigGetResp, err error) {
	if req == nil || req.Key == "" {
		return nil, errs.New(errs.CodeBadRequest, "配置键不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.ConfigGet(l.ctx, &iamclient.ConfigGetRequest{Key: req.Key})
	if err != nil {
		return nil, errs.WrapGRPCError("查询配置失败", err)
	}

	return &types.ConfigGetResp{Value: rpcResp.Value}, nil
}
