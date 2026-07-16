package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleUpdateLogic {
	return &UserRoleUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleUpdateLogic) UserRoleUpdate(in *iam.UserRoleUpdateRequest) (*iam.Empty, error) {
	if in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户ID不能为空"))
	}

	if err := l.svcCtx.Domain.IAM.RBAC.UpdateUserRoles(l.ctx, in.UserId, in.RoleIds); err != nil {
		return nil, toGRPCStatus(err)
	}

	cache := l.svcCtx.Repository.BusinessCache
	userID := in.UserId
	go func() {
		if err := cache.DeleteUserPermissions(context.Background(), userID); err != nil {
			l.Errorf("清除用户权限缓存失败: userId=%d, error=%v", userID, err)
		}
		if err := cache.DeleteUserMenuTree(context.Background(), userID); err != nil {
			l.Errorf("清除用户菜单树缓存失败: userId=%d, error=%v", userID, err)
		}
	}()

	return &iam.Empty{}, nil
}
