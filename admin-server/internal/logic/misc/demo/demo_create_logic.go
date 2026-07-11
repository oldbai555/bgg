// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package demo

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"postapocgame/admin-server/internal/model/misc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCreateLogic {
	return &DemoCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoCreateLogic) DemoCreate(req *types.DemoCreateReq) error {
	if req == nil || req.Name == "" {
		return errs.New(errs.CodeBadRequest, "描述不能为空")
	}

	status := req.Status
	if status == 0 {
		status = 1
	}

	demo := misc.Demo{
		Name:   req.Name,
		Status: status,
	}

	if err := l.svcCtx.Domain.Misc.Demo.Create(l.ctx, &demo); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建演示功能失败", err)
	}
	return nil
}
