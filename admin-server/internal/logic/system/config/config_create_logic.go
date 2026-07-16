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

type ConfigCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigCreateLogic {
	return &ConfigCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigCreateLogic) ConfigCreate(req *types.ConfigCreateReq) error {
	if req == nil || req.Group == "" || req.Key == "" {
		return errs.New(errs.CodeBadRequest, "配置分组和键不能为空")
	}

	_, err := l.svcCtx.IamRPC.ConfigCreate(l.ctx, &iamclient.ConfigCreateRequest{
		Group:       req.Group,
		Key:         req.Key,
		Value:       req.Value,
		ConfigType:  req.ConfigType,
		Description: req.Description,
	})
	if err != nil {
		return errs.WrapGRPCError("创建配置失败", err)
	}
	return nil
}
