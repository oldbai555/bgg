package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/initdata"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPermissionDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionDeleteLogic {
	return &PermissionDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PermissionDeleteLogic) PermissionDelete(in *iam.PermissionDeleteRequest) (*iam.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "权限ID不能为空"))
	}
	if initdata.IsInitPermissionID(in.Id) {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "初始化数据不可删除"))
	}

	if err := l.svcCtx.Domain.IAM.Permission.DeleteByID(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "删除权限失败", err))
	}
	return &iam.Empty{}, nil
}
