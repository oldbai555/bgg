package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateDaySaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTemplateDaySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateDaySaveLogic {
	return &FitnessTemplateDaySaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTemplateDaySaveLogic) FitnessTemplateDaySave(in *iam.FitnessTemplateDaySaveRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.SaveTemplateDay(l.ctx, in.TemplateId, in.Weekday, in.PlanDate, fitnessDayContentInput(in.Content)); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
