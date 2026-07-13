package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfigGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigGetLogic {
	return &ConfigGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfigGetLogic) ConfigGet(in *iam.ConfigGetRequest) (*iam.ConfigGetResponse, error) {
	if in == nil || in.Key == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "配置键不能为空"))
	}

	cache := l.svcCtx.Repository.BusinessCache
	var cachedValue string
	if err := cache.GetConfigKey(l.ctx, in.Key, &cachedValue); err == nil {
		return &iam.ConfigGetResponse{Value: cachedValue}, nil
	}

	cfg, err := l.svcCtx.Domain.System.Config.FindByKey(l.ctx, in.Key)
	if err != nil {
		if isErrNotFound(err) {
			return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "配置不存在"))
		}
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询配置失败", err))
	}

	value := ""
	if cfg.Value.Valid {
		value = cfg.Value.String
	}

	key := in.Key
	go func() {
		if err := cache.SetConfigKey(context.Background(), key, value); err != nil {
			l.Errorf("设置配置缓存失败: key=%s, error=%v", key, err)
		}
	}()

	return &iam.ConfigGetResponse{Value: value}, nil
}
