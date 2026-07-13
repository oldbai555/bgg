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

type DepartmentTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDepartmentTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DepartmentTreeLogic {
	return &DepartmentTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DepartmentTreeLogic) DepartmentTree() (resp *types.DepartmentTreeResp, err error) {
	rpcResp, err := l.svcCtx.IamRPC.DepartmentTree(l.ctx, &iamclient.Empty{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询部门列表失败", err)
	}

	return &types.DepartmentTreeResp{List: convertDepartmentItems(rpcResp.List)}, nil
}

func convertDepartmentItems(items []*iamclient.DepartmentItem) []types.DepartmentItem {
	result := make([]types.DepartmentItem, 0, len(items))
	for _, item := range items {
		result = append(result, types.DepartmentItem{
			Id:       item.Id,
			ParentId: item.ParentId,
			Name:     item.Name,
			OrderNum: item.OrderNum,
			Status:   item.Status,
			Children: convertDepartmentItems(item.Children),
		})
	}
	return result
}
