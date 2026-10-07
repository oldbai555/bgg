package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTemplateCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateCreateLogic {
	return &FitnessTemplateCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTemplateCreateLogic) FitnessTemplateCreate(in *iam.FitnessTemplateSaveRequest) (*iam.Empty, error) {
	if err := l.svcCtx.Domain.Fitness.CreateTemplate(l.ctx, fitnessTemplateModel(in)); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
