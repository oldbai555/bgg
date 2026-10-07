package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	fitnessdomain "postapocgame/admin-server/services/iam/internal/domain/fitness"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyDayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMyDayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyDayLogic {
	return &FitnessMyDayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMyDayLogic) FitnessMyDay(in *iam.FitnessMyDayRequest) (*iam.FitnessMyDayResponse, error) {
	plan, err := l.svcCtx.Domain.Fitness.DayPlan(l.ctx, in.UserId, in.Date)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	suggestion, err := l.svcCtx.Domain.Fitness.CurrentSuggestion(l.ctx, in.UserId)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessMyDayResponse{
		Date:       plan.Date,
		Today:      plan.Today,
		Weekday:    fitnessdomain.Weekday(plan.Date),
		CanCheckin: plan.CanCheckin,
		HasPlan:    plan.Day != nil,
		IsOverride: plan.IsOverride,
		Content:    fitnessDayContent(plan.Day),
		Checkin:    fitnessCheckinItem(plan.Checkin),
		Suggestion: fitnessSuggestionItem(suggestion),
	}
	if plan.Template != nil {
		resp.TemplateId, resp.TemplateName = plan.Template.Id, plan.Template.Name
	}
	for _, t := range plan.Tips {
		resp.Tips = append(resp.Tips, fitnessTipItem(t))
	}
	return resp, nil
}
