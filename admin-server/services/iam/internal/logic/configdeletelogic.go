package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigDeleteLogic {
	return &ConfigDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigDeleteLogic) ConfigDelete(in *iam.ConfigDeleteRequest) (*iam.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "配置ID不能为空"))
	}

	cfg, err := l.svcCtx.Domain.System.Config.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询配置失败", err))
	}

	if err := l.svcCtx.Domain.System.Config.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除配置失败", err))
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
