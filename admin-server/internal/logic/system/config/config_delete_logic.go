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

type ConfigDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigDeleteLogic {
	return &ConfigDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigDeleteLogic) ConfigDelete(req *types.ConfigDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "配置ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.ConfigDelete(l.ctx, &iamclient.ConfigDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除配置失败", err)
	}
	return nil
}
