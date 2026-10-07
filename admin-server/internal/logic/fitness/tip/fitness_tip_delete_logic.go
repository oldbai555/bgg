// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tip

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTipDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTipDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipDeleteLogic {
	return &FitnessTipDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTipDeleteLogic) FitnessTipDelete(req *types.FitnessIdReq) error {
	if _, err := l.svcCtx.IamRPC.FitnessTipDelete(l.ctx, &iamclient.FitnessIdRequest{Id: req.Id}); err != nil {
		return errs.WrapGRPCError("删除通用提示失败", err)
	}
	return nil
}
