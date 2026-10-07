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

type FitnessTipCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTipCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipCreateLogic {
	return &FitnessTipCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTipCreateLogic) FitnessTipCreate(req *types.FitnessTipSaveReq) error {
	_, err := l.svcCtx.IamRPC.FitnessTipCreate(l.ctx, &iamclient.FitnessTipSaveRequest{
		Id:      req.Id,
		Code:    req.Code,
		Title:   req.Title,
		Content: req.Content,
		Sort:    req.Sort,
		Status:  req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("创建通用提示失败", err)
	}
	return nil
}
