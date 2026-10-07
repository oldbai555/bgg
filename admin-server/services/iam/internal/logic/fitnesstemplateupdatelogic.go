package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFitnessTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateUpdateLogic {
	return &FitnessTemplateUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FitnessTemplateUpdateLogic) FitnessTemplateUpdate(in *iam.FitnessTemplateSaveRequest) (*iam.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "id不能为空"))
	}
	if err := l.svcCtx.Domain.Fitness.UpdateTemplate(l.ctx, fitnessTemplateModel(in)); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
