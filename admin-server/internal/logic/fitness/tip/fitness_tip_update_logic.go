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

type FitnessTipUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTipUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipUpdateLogic {
	return &FitnessTipUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTipUpdateLogic) FitnessTipUpdate(req *types.FitnessTipSaveReq) error {
	_, err := l.svcCtx.IamRPC.FitnessTipUpdate(l.ctx, &iamclient.FitnessTipSaveRequest{
		Id:      req.Id,
		Code:    req.Code,
		Title:   req.Title,
		Content: req.Content,
		Sort:    req.Sort,
		Status:  req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("更新通用提示失败", err)
	}
	return nil
}
