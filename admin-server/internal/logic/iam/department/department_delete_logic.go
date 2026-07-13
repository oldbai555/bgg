// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package department

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDepartmentDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentDeleteLogic {
	return &DepartmentDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentDeleteLogic) DepartmentDelete(req *types.DepartmentDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "部门ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.DepartmentDelete(l.ctx, &iamclient.DepartmentDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除部门失败", err)
	}
	return nil
}
