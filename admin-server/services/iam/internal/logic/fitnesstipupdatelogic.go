package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTipUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTipUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipUpdateLogic {
	return &FitnessTipUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTipUpdateLogic) FitnessTipUpdate(in *iam.FitnessTipSaveRequest) (*iam.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "id不能为空"))
	}
	if err := l.svcCtx.Domain.Fitness.UpdateTip(l.ctx, fitnessTipModel(in)); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
