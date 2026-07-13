// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package permission_api

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PermissionApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPermissionApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PermissionApiListLogic {
	return &PermissionApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PermissionApiListLogic) PermissionApiList(req *types.PermissionApiListReq) (resp *types.PermissionApiListResp, err error) {
	if req.PermissionId == 0 {
		return nil, errs.New(errs.CodeBadRequest, "权限ID不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.PermissionApiList(l.ctx, &iamclient.PermissionApiListRequest{PermissionId: req.PermissionId})
	if err != nil {
		return nil, errs.WrapGRPCError("查询权限接口失败", err)
	}

	return &types.PermissionApiListResp{ApiIds: rpcResp.ApiIds}, nil
}
