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

type ConfigUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigUpdateLogic {
	return &ConfigUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigUpdateLogic) ConfigUpdate(req *types.ConfigUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "配置ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.ConfigUpdate(l.ctx, &iamclient.ConfigUpdateRequest{
		Id:          req.Id,
		Value:       req.Value,
		Description: req.Description,
	})
	if err != nil {
		return errs.WrapGRPCError("更新配置失败", err)
	}
	return nil
}
