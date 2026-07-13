// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleListLogic) RoleList(req *types.RoleListReq) (resp *types.RoleListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.RoleList(l.ctx, &iamclient.RoleListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询角色列表失败", err)
	}

	items := make([]types.RoleItem, 0, len(rpcResp.List))
	for _, r := range rpcResp.List {
		items = append(items, types.RoleItem{
			Id:          r.Id,
			Name:        r.Name,
			Code:        r.Code,
			Description: r.Description,
			Status:      r.Status,
		})
	}

	return &types.RoleListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
