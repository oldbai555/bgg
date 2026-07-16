package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 权限校验（GetUserProfile 等见 pkg/iamcallback.IamCallback，同一进程内另注册）
func (l *CheckPermissionLogic) CheckPermission(in *iam.CheckPermissionRequest) (*iam.CheckPermissionResponse, error) {
	allowed, err := l.svcCtx.Domain.IAM.PermissionResolver.CanAccess(l.ctx, in.UserId, in.Method, in.Path)
	if err != nil {
		return &iam.CheckPermissionResponse{Allowed: false, Reason: err.Error()}, nil
	}
	return &iam.CheckPermissionResponse{Allowed: allowed}, nil
}
