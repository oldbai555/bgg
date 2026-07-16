package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProfileLogic) Profile(in *iam.ProfileRequest) (*iam.ProfileResponse, error) {
	if in == nil || in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}

	userInfo, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, in.UserId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "获取用户信息失败", err))
	}

	cache := l.svcCtx.Repository.BusinessCache
	codes, err := cache.GetUserPermissions(l.ctx, in.UserId)
	if err != nil {
		roleIDs, err := l.svcCtx.Domain.IAM.UserRole.ListRoleIDsByUserID(l.ctx, in.UserId)
		if err != nil {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "获取用户角色失败", err))
		}

		perms, err := l.svcCtx.Domain.IAM.Permission.ListByRoleIDs(l.ctx, roleIDs)
		if err != nil {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "获取用户权限失败", err))
		}

		codes = make([]string, 0, len(perms))
		seen := make(map[string]struct{}, len(perms))
		for _, p := range perms {
			if _, ok := seen[p.Code]; ok {
				continue
			}
			seen[p.Code] = struct{}{}
			codes = append(codes, p.Code)
		}

		userID := in.UserId
		go func() {
			if err := cache.SetUserPermissions(context.Background(), userID, codes); err != nil {
				l.Errorf("设置用户权限缓存失败: userId=%d, error=%v", userID, err)
			}
		}()
	}

	return &iam.ProfileResponse{
		Id:          in.UserId,
		Username:    userInfo.Username,
		Nickname:    userInfo.Nickname,
		Avatar:      userInfo.Avatar,
		Signature:   userInfo.Signature,
		Permissions: codes,
	}, nil
}
