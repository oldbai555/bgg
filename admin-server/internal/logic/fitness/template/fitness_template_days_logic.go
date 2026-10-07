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

type FitnessTemplateDaysLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTemplateDaysLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateDaysLogic {
	return &FitnessTemplateDaysLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTemplateDaysLogic) FitnessTemplateDays(req *types.FitnessTemplateDaysReq) (resp *types.FitnessTemplateDaysResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.FitnessTemplateDays(l.ctx, &iamclient.FitnessTemplateDaysRequest{TemplateId: req.TemplateId})
	if err != nil {
		return nil, errs.WrapGRPCError("查询模板日计划失败", err)
	}
	return &types.FitnessTemplateDaysResp{
		Weekly:    fitnessconv.TemplateDayItems(rpcResp.Weekly),
		Overrides: fitnessconv.TemplateDayItems(rpcResp.Overrides),
	}, nil
}
