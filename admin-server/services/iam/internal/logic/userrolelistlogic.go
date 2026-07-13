package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleListLogic {
	return &UserRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleListLogic) UserRoleList(in *iam.UserRoleListRequest) (*iam.UserRoleListResponse, error) {
	if in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户ID不能为空"))
	}

	if _, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, in.UserId); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadRequest, "用户不存在", err))
	}

	roleIDs, err := l.svcCtx.Domain.IAM.UserRole.ListRoleIDsByUserID(l.ctx, in.UserId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询用户角色失败", err))
	}

	return &iam.UserRoleListResponse{RoleIds: roleIDs}, nil
}
