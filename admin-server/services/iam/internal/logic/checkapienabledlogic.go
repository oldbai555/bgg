package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckApiEnabledLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckApiEnabledLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckApiEnabledLogic {
	return &CheckApiEnabledLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Gateway 基础设施
func (l *CheckApiEnabledLogic) CheckApiEnabled(in *iam.CheckApiEnabledRequest) (*iam.CheckApiEnabledResponse, error) {
	api, err := l.svcCtx.Domain.IAM.Api.FindByMethodAndPath(l.ctx, in.Method, in.Path)
	if err != nil {
		return &iam.CheckApiEnabledResponse{Exists: false, Enabled: false}, nil
	}
	return &iam.CheckApiEnabledResponse{Exists: true, Enabled: api.Status == 1}, nil
}
