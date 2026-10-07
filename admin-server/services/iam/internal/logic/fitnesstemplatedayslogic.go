package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateDaysLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTemplateDaysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateDaysLogic {
	return &FitnessTemplateDaysLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTemplateDaysLogic) FitnessTemplateDays(in *iam.FitnessTemplateDaysRequest) (*iam.FitnessTemplateDaysResponse, error) {
	weekly, overrides, err := l.svcCtx.Domain.Fitness.TemplateDays(l.ctx, in.TemplateId)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessTemplateDaysResponse{}
	for _, d := range weekly {
		resp.Weekly = append(resp.Weekly, fitnessTemplateDayItem(d))
	}
	for _, d := range overrides {
		resp.Overrides = append(resp.Overrides, fitnessTemplateDayItem(d))
	}
	return resp, nil
}
