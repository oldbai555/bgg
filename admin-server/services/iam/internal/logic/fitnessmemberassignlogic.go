package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessMemberAssignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessMemberAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMemberAssignLogic {
	return &FitnessMemberAssignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessMemberAssignLogic) FitnessMemberAssign(in *iam.FitnessMemberAssignRequest) (*iam.FitnessTemplateSwitchResponse, error) {
	effective, err := l.svcCtx.Domain.Fitness.AssignTemplate(l.ctx, in.UserId, in.TemplateId, in.Reason, in.OperatorId)
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.FitnessTemplateSwitchResponse{EffectiveDate: effective}, nil
}
