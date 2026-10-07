package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyBodyRecordSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMyBodyRecordSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyBodyRecordSaveLogic {
	return &FitnessMyBodyRecordSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMyBodyRecordSaveLogic) FitnessMyBodyRecordSave(in *iam.FitnessMyBodyRecordSaveRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.SaveBodyRecord(l.ctx, in.UserId, in.Date, in.WeightKg, in.WaistCm); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
