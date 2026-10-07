package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateDayDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTemplateDayDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateDayDeleteLogic {
	return &FitnessTemplateDayDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTemplateDayDeleteLogic) FitnessTemplateDayDelete(in *iam.FitnessIdRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.DeleteTemplateDay(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
