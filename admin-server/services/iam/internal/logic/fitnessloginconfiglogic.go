package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessLoginConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessLoginConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessLoginConfigLogic {
	return &FitnessLoginConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Fitness 手机端：只操作 user_id 本人的数据
func (l *FitnessLoginConfigLogic) FitnessLoginConfig(in *iam.Empty) (*iam.FitnessLoginConfigResponse, error) {
	return &iam.FitnessLoginConfigResponse{AppId: l.svcCtx.Config.Feishu.AppId}, nil
}
