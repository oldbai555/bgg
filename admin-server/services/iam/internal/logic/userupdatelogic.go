package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserUpdateLogic) UserUpdate(in *iam.UserUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户ID不能为空"))
	}

	user, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, in.Id)
	if err != nil {
		if isErrNotFound(err) {
			return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "用户不存在"))
		}
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询用户失败", err))
	}

	if in.Username != "" {
		existing, err := l.svcCtx.Domain.IAM.User.FindByUsername(l.ctx, in.Username)
		if err == nil && existing.Id != in.Id {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户名已被使用"))
		}
		if err != nil && !isErrNotFound(err) {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询用户名失败", err))
		}
		user.Username = in.Username
	}

	if in.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "密码加密失败", err))
		}
		user.PasswordHash = string(hash)
	}

	if in.DepartmentId != 0 {
		user.DepartmentId = in.DepartmentId
	}

	if in.Nickname != "" {
		user.Nickname = in.Nickname
	}
	if in.Avatar != "" {
		user.Avatar = in.Avatar
	}
	if in.Signature != "" {
		user.Signature = in.Signature
	}

	if in.Status == 0 || in.Status == 1 {
		user.Status = in.Status
	}

	if err := l.svcCtx.Domain.IAM.User.Update(l.ctx, user); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新用户失败", err))
	}
	return &iam.Empty{}, nil
}
