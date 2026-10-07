package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTipCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTipCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipCreateLogic {
	return &FitnessTipCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTipCreateLogic) FitnessTipCreate(in *iam.FitnessTipSaveRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.CreateTip(l.ctx, fitnessTipModel(in)); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
