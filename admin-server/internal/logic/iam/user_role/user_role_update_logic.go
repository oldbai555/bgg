// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user_role

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleUpdateLogic {
	return &UserRoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserRoleUpdateLogic) UserRoleUpdate(req *types.UserRoleUpdateReq) error {
	if req.UserId == 0 {
		return errs.New(errs.CodeBadRequest, "用户ID不能为空")
	}

	if err := l.svcCtx.Domain.IAM.RBAC.UpdateUserRoles(l.ctx, req.UserId, req.RoleIds); err != nil {
		return err
	}

	// 清除该用户的权限和菜单树缓存
	cache := l.svcCtx.Repository.BusinessCache
	go func() {
		if err := cache.DeleteUserPermissions(context.Background(), req.UserId); err != nil {
			l.Errorf("清除用户权限缓存失败: userId=%d, error=%v", req.UserId, err)
		}
		if err := cache.DeleteUserMenuTree(context.Background(), req.UserId); err != nil {
			l.Errorf("清除用户菜单树缓存失败: userId=%d, error=%v", req.UserId, err)
		}
	}()

	return nil
}
