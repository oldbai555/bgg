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

type FitnessTemplateDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTemplateDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateDeleteLogic {
	return &FitnessTemplateDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTemplateDeleteLogic) FitnessTemplateDelete(req *types.FitnessIdReq) error {
	if _, err := l.svcCtx.IamRPC.FitnessTemplateDelete(l.ctx, &iamclient.FitnessIdRequest{Id: req.Id}); err != nil {
		return errs.WrapGRPCError("删除计划模板失败", err)
	}
	return nil
}
