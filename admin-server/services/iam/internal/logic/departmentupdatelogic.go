package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDepartmentUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentUpdateLogic {
	return &DepartmentUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DepartmentUpdateLogic) DepartmentUpdate(in *iam.DepartmentUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "部门ID不能为空"))
	}

	dept, err := l.svcCtx.Domain.IAM.Department.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询部门失败", err))
	}

	dept.ParentId = in.ParentId
	dept.Name = in.Name
	if in.OrderNum >= 0 {
		dept.OrderNum = in.OrderNum
	}
	if in.Status == 0 || in.Status == 1 {
		dept.Status = in.Status
	}

	if err := l.svcCtx.Domain.IAM.Department.Update(l.ctx, dept); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新部门失败", err))
	}
	return &iam.Empty{}, nil
}
