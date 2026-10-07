// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package template

import (
	"context"

	"postapocgame/admin-server/internal/logic/fitness/fitnessconv"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTemplateDaySaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTemplateDaySaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateDaySaveLogic {
	return &FitnessTemplateDaySaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTemplateDaySaveLogic) FitnessTemplateDaySave(req *types.FitnessTemplateDaySaveReq) error {
	_, err := l.svcCtx.IamRPC.FitnessTemplateDaySave(l.ctx, &iamclient.FitnessTemplateDaySaveRequest{
		TemplateId: req.TemplateId,
		Weekday:    req.Weekday,
		PlanDate:   req.PlanDate,
		Content:    fitnessconv.DayContentToRPC(req.Content),
	})
	if err != nil {
		return errs.WrapGRPCError("保存模板日计划失败", err)
	}
	return nil
}
