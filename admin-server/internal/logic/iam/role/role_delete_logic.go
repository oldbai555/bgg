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

type RoleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleDeleteLogic {
	return &RoleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleDeleteLogic) RoleDelete(req *types.RoleDeleteReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "角色ID不能为空")
	}

	_, err := l.svcCtx.IamRPC.RoleDelete(l.ctx, &iamclient.RoleDeleteRequest{Id: req.Id})
	if err != nil {
		return errs.WrapGRPCError("删除角色失败", err)
	}
	return nil
}
