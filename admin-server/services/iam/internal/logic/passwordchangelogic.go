package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type PasswordChangeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPasswordChangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordChangeLogic {
	return &PasswordChangeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PasswordChangeLogic) PasswordChange(in *iam.PasswordChangeRequest) (*iam.Empty, error) {
	if in == nil || in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}
	if in.OldPassword == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "原密码不能为空"))
	}
	if in.NewPassword == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "新密码不能为空"))
	}
	if len(in.NewPassword) < 6 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "新密码长度不能少于6位"))
	}

	userInfo, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, in.UserId)
	if err != nil {
		if isErrNotFound(err) {
			return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "用户不存在"))
		}
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "获取用户信息失败", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.PasswordHash), []byte(in.OldPassword)); err != nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "原密码错误"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "加密密码失败", err))
	}

	userInfo.PasswordHash = string(hashedPassword)
	userInfo.UpdatedAt = time.Now().Unix()

	if err := l.svcCtx.Domain.IAM.User.Update(l.ctx, userInfo); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新密码失败", err))
	}

	return &iam.Empty{}, nil
}
