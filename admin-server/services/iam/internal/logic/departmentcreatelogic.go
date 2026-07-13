package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDepartmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentCreateLogic {
	return &DepartmentCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DepartmentCreateLogic) DepartmentCreate(in *iam.DepartmentCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Name == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "部门名称不能为空"))
	}

	dept := iammodel.AdminDepartment{
		ParentId: in.ParentId,
		Name:     in.Name,
		OrderNum: in.OrderNum,
		Status:   in.Status,
	}

	if err := l.svcCtx.Domain.IAM.Department.Create(l.ctx, &dept); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "创建部门失败", err))
	}
	return &iam.Empty{}, nil
}
