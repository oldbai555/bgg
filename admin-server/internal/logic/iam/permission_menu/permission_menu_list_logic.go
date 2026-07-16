// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package permission_menu

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionMenuListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionMenuListLogic {
	return &PermissionMenuListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionMenuListLogic) PermissionMenuList(req *types.PermissionMenuListReq) (resp *types.PermissionMenuListResp, err error) {
	if req.PermissionId == 0 {
		return nil, errs.New(errs.CodeBadRequest, "权限ID不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.PermissionMenuList(l.ctx, &iamclient.PermissionMenuListRequest{PermissionId: req.PermissionId})
	if err != nil {
		return nil, errs.WrapGRPCError("查询权限菜单失败", err)
	}

	return &types.PermissionMenuListResp{MenuIds: rpcResp.MenuIds}, nil
}
