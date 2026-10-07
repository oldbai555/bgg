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

type FitnessTemplateUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTemplateUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateUpdateLogic {
	return &FitnessTemplateUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTemplateUpdateLogic) FitnessTemplateUpdate(req *types.FitnessTemplateUpdateReq) error {
	_, err := l.svcCtx.IamRPC.FitnessTemplateUpdate(l.ctx, &iamclient.FitnessTemplateSaveRequest{
		Id:           req.Id,
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
		return errs.WrapGRPCError("更新计划模板失败", err)
	}
	return nil
}
