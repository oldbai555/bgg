package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTipDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTipDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipDeleteLogic {
	return &FitnessTipDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTipDeleteLogic) FitnessTipDelete(in *iam.FitnessIdRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.DeleteTip(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
