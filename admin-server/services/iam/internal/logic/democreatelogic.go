package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	miscmodel "postapocgame/admin-server/services/iam/internal/model/misc"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDemoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoCreateLogic {
	return &DemoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Demo / DailyShortSentence
func (l *DemoCreateLogic) DemoCreate(in *iam.DemoCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "描述不能为空"))
	}

	status := in.Status
	if status == 0 {
		status = 1
	}

	demo := miscmodel.Demo{
		Name:   in.Name,
		Status: status,
	}

	if err := l.svcCtx.Domain.Misc.Demo.Create(l.ctx, &demo); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建演示功能失败", err))
	}
	return &iam.Empty{}, nil
}
