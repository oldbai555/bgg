package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigUpdateLogic {
	return &ConfigUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigUpdateLogic) ConfigUpdate(in *iam.ConfigUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "配置ID不能为空"))
	}

	cfg, err := l.svcCtx.Domain.System.Config.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询配置失败", err))
	}

	if in.Value != "" {
		cfg.Value = sql.NullString{String: in.Value, Valid: true}
	}
	if in.Description != "" {
		cfg.Description = sql.NullString{String: in.Description, Valid: true}
	}

	if err := l.svcCtx.Domain.System.Config.Update(l.ctx, cfg); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新配置失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	key := cfg.Key
	go func() {
		if err := cache.DeleteConfigKey(context.Background(), key); err != nil {
			l.Errorf("清除配置缓存失败: key=%s, error=%v", key, err)
		}
	}()

	return &iam.Empty{}, nil
}
