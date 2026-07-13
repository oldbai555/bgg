package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDemoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoUpdateLogic {
	return &DemoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DemoUpdateLogic) DemoUpdate(in *iam.DemoUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	demo, err := l.svcCtx.Domain.Misc.Demo.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "演示功能不存在", err))
	}

	if in.Name != "" {
		demo.Name = in.Name
	}
	if in.Status == 0 || in.Status == 1 {
		demo.Status = in.Status
	}

	if err := l.svcCtx.Domain.Misc.Demo.Update(l.ctx, demo); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新演示功能失败", err))
	}
	return &iam.Empty{}, nil
}
