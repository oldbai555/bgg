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

type DemoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoDeleteLogic {
	return &DemoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoDeleteLogic) DemoDelete(req *types.DemoDeleteReq) error {
	if req == nil || req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	_, err := l.svcCtx.IamRPC.DemoDelete(l.ctx, &iamclient.DemoDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除演示功能失败", err)
	}
	return nil
}
