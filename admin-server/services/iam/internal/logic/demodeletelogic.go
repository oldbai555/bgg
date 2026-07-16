package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDemoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoDeleteLogic {
	return &DemoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DemoDeleteLogic) DemoDelete(in *iam.DemoDeleteRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	if err := l.svcCtx.Domain.Misc.Demo.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除演示功能失败", err))
	}
	return &iam.Empty{}, nil
}
