// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package demo

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoUpdateLogic {
	return &DemoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoUpdateLogic) DemoUpdate(req *types.DemoUpdateReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	_, err := l.svcCtx.IamRPC.DemoUpdate(l.ctx, &iamclient.DemoUpdateRequest{
		Id:     req.Id,
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return errs.WrapGRPCError("更新演示功能失败", err)
	}
	return nil
}
