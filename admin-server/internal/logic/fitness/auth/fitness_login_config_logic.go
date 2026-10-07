// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessLoginConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessLoginConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessLoginConfigLogic {
	return &FitnessLoginConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessLoginConfigLogic) FitnessLoginConfig() (resp *types.FitnessLoginConfigResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.FitnessLoginConfig(l.ctx, &iamclient.Empty{})
	if err != nil {
		return nil, errs.WrapGRPCError("获取登录配置失败", err)
	}
	return &types.FitnessLoginConfigResp{AppId: rpcResp.AppId}, nil
}
