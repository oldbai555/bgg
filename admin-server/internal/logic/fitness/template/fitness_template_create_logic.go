// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package template

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTemplateCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateCreateLogic {
	return &FitnessTemplateCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTemplateCreateLogic) FitnessTemplateCreate(req *types.FitnessTemplateCreateReq) error {
	_, err := l.svcCtx.IamRPC.FitnessTemplateCreate(l.ctx, &iamclient.FitnessTemplateSaveRequest{
		Code:         req.Code,
		Name:         req.Name,
		Level:        req.Level,
		Stage:        req.Stage,
		Summary:      req.Summary,
		DailyKcal:    req.DailyKcal,
		DailyProtein: req.DailyProtein,
		Sort:         req.Sort,
		Status:       req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("创建计划模板失败", err)
	}
	return nil
}
