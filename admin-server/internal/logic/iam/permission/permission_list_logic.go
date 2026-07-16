// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package permission

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionListLogic {
	return &PermissionListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionListLogic) PermissionList(req *types.PermissionListReq) (resp *types.PermissionListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.PermissionList(l.ctx, &iamclient.PermissionListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询权限列表失败", err)
	}

	items := make([]types.PermissionItem, 0, len(rpcResp.List))
	for _, p := range rpcResp.List {
		items = append(items, types.PermissionItem{
			Id:          p.Id,
			Name:        p.Name,
			Code:        p.Code,
			Description: p.Description,
		})
	}

	return &types.PermissionListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
