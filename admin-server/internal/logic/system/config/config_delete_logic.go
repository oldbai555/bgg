// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

	// 先查询配置，获取 key
	config, err := l.svcCtx.Domain.System.Config.FindByID(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, "查询配置失败", err)
	}

	if err := l.svcCtx.Domain.System.Config.DeleteByID(l.ctx, req.Id); err != nil {
		return errs.Wrap(errs.CodeInternalError, "删除配置失败", err)
	}

	// 清除配置缓存
	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteConfigKey(context.Background(), config.Key); err != nil {
			l.Errorf("清除配置缓存失败: key=%s, error=%v", config.Key, err)
		}
	}()

	return nil
}
