package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	iamdomain "postapocgame/admin-server/services/iam/internal/domain/iam"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// User / Role / Permission / Menu / Department / Api
func (l *UserCreateLogic) UserCreate(in *iam.UserCreateRequest) (*iam.Empty, error) {
	if in == nil || in.Username == "" || in.Password == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户名和密码不能为空"))
	}

	_, err := l.svcCtx.Domain.IAM.UserService.CreateUser(l.ctx, iamdomain.CreateUserInput{
		Username:     in.Username,
		Nickname:     in.Nickname,
		Password:     in.Password,
		Avatar:       in.Avatar,
		Signature:    in.Signature,
		DepartmentId: in.DepartmentId,
		Status:       in.Status,
	})
	if err != nil {
		return nil, toGRPCStatus(err)
	}
	return &iam.Empty{}, nil
}
