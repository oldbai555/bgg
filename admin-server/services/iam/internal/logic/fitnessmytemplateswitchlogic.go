package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/consts"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMyTemplateSwitchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMyTemplateSwitchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMyTemplateSwitchLogic {
	return &FitnessMyTemplateSwitchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMyTemplateSwitchLogic) FitnessMyTemplateSwitch(in *iam.FitnessMyTemplateSwitchRequest) (*iam.FitnessTemplateSwitchResponse, error) {
	if _, err := l.svcCtx.Domain.Fitness.EnsureMember(l.ctx, in.UserId); err != nil {
		return nil, toGRPCStatus(err)
	}
	effective, err := l.svcCtx.Domain.Fitness.SwitchTemplate(l.ctx, in.UserId, in.TemplateId, consts.FitnessSwitchSourceManual, consts.FitnessSwitchReasonManual, in.UserId)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.FitnessTemplateSwitchResponse{EffectiveDate: effective}, nil
}
