package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMyProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyProfileLogic {
	return &FitnessMyProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMyProfileLogic) FitnessMyProfile(in *iam.FitnessMyProfileRequest) (*iam.FitnessMyProfileResponse, error) {
	p, err := l.svcCtx.Domain.Fitness.MyProfile(l.ctx, in.UserId)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	resp := &iam.FitnessMyProfileResponse{
		Today:              p.Today,
		Nickname:           p.Nickname,
		Avatar:             p.Avatar,
		HeightCm:           p.Profile.HeightCm,
		TargetWeightKg:     p.Profile.TargetWeightKg,
		TrainingPreference: p.Profile.TrainingPreference,
		CurrentTemplate:    fitnessTemplateBrief(p.Current),
		Suggestion:         fitnessSuggestionItem(p.Suggestion),
	}
	if p.Member != nil {
		resp.LatestWeightKg, resp.StreakDays, resp.WeekRate = p.Member.LatestWeight, int64(p.Member.StreakDays), p.Member.WeekRate
	}
	if p.Pending != nil {
		resp.PendingSwitch = &iam.FitnessPendingSwitch{TemplateId: p.Pending.ToTemplateId, EffectiveDate: p.Pending.EffectiveDate}
		if p.PendingTpl != nil {
			resp.PendingSwitch.TemplateName = p.PendingTpl.Name
		}
	}
	for i := range p.Templates {
		resp.Templates = append(resp.Templates, fitnessTemplateBrief(&p.Templates[i]))
	}
	return resp, nil
}
