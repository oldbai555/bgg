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

type FitnessTemplateDayDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTemplateDayDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateDayDeleteLogic {
	return &FitnessTemplateDayDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTemplateDayDeleteLogic) FitnessTemplateDayDelete(req *types.FitnessIdReq) error {
	if _, err := l.svcCtx.IamRPC.FitnessTemplateDayDelete(l.ctx, &iamclient.FitnessIdRequest{Id: req.Id}); err != nil {
		return errs.WrapGRPCError("删除日期调整失败", err)
	}
	return nil
}
