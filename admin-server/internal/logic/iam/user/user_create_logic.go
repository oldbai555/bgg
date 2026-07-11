// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	iamdomain "postapocgame/admin-server/internal/domain/iam"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCreateLogic) UserCreate(req *types.UserCreateReq) error {
	if req == nil || req.Username == "" || req.Password == "" {
		return errs.New(errs.CodeBadRequest, "用户名和密码不能为空")
	}

	_, err := l.svcCtx.Domain.IAM.UserService.CreateUser(l.ctx, iamdomain.CreateUserInput{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Password:     req.Password,
		Avatar:       req.Avatar,
		Signature:    req.Signature,
		DepartmentId: req.DepartmentId,
		Status:       req.Status,
	})
	return err
}
