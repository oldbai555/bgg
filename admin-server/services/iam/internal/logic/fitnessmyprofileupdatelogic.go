package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyProfileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMyProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyProfileUpdateLogic {
	return &FitnessMyProfileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMyProfileUpdateLogic) FitnessMyProfileUpdate(in *iam.FitnessMyProfileUpdateRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.UpdateProfile(l.ctx, in.UserId, in.HeightCm, in.TargetWeightKg, in.TrainingPreference); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
