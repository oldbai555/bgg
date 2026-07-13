package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DepartmentTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDepartmentTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentTreeLogic {
	return &DepartmentTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DepartmentTreeLogic) DepartmentTree(in *iam.Empty) (*iam.DepartmentTreeResponse, error) {
	list, err := l.svcCtx.Domain.IAM.Department.ListAll(l.ctx)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询部门列表失败", err))
	}

	nodeMap := make(map[uint64]*iam.DepartmentItem, len(list))
	var roots []*iam.DepartmentItem

	for _, d := range list {
		item := &iam.DepartmentItem{
			Id:       d.Id,
			ParentId: d.ParentId,
			Name:     d.Name,
			OrderNum: d.OrderNum,
			Status:   d.Status,
		}
		nodeMap[d.Id] = item
	}

	for _, item := range nodeMap {
		if item.ParentId == 0 {
			roots = append(roots, item)
			continue
		}
		if parent, ok := nodeMap[item.ParentId]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			roots = append(roots, item)
		}
	}

	return &iam.DepartmentTreeResponse{List: roots}, nil
}
