// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package department

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"postapocgame/admin-server/internal/model/iam"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDepartmentCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentCreateLogic {
	return &DepartmentCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentCreateLogic) DepartmentCreate(req *types.DepartmentCreateReq) error {
	if req == nil || req.Name == "" {
		return errs.New(errs.CodeBadRequest, "部门名称不能为空")
	}

	dept := iam.AdminDepartment{
		ParentId: req.ParentId,
		Name:     req.Name,
		OrderNum: req.OrderNum,
		Status:   req.Status,
	}

	if err := l.svcCtx.Domain.IAM.Department.Create(l.ctx, &dept); err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建部门失败", err)
	}
	return nil
}
