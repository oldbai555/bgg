package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	systemmodel "postapocgame/admin-server/services/iam/internal/model/system"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigCreateLogic {
	return &ConfigCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Config / Dict
func (l *ConfigCreateLogic) ConfigCreate(in *iam.ConfigCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Group == "" || in.Key == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "配置分组和键不能为空"))
	}

	if _, err := l.svcCtx.Domain.System.Config.FindByKey(l.ctx, in.Key); err == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "配置键已存在"))
	}

	configType := in.ConfigType
	if configType == "" {
		configType = "string"
	}

	cfg := systemmodel.AdminConfig{
		Group:       in.Group,
		Key:         in.Key,
		Value:       sql.NullString{String: in.Value, Valid: in.Value != ""},
		Type:        configType,
		Description: sql.NullString{String: in.Description, Valid: in.Description != ""},
	}

	if err := l.svcCtx.Domain.System.Config.Create(l.ctx, &cfg); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建配置失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	key, value := in.Key, in.Value
	go func() {
		if value != "" {
			if err := cache.SetConfigKey(context.Background(), key, value); err != nil {
				l.Errorf("设置配置缓存失败: key=%s, error=%v", key, err)
			}
		} else {
			if err := cache.DeleteConfigKey(context.Background(), key); err != nil {
				l.Errorf("清除配置缓存失败: key=%s, error=%v", key, err)
			}
		}
	}()

	return &iam.Empty{}, nil
}
