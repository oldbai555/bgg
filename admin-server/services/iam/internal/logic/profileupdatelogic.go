package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProfileUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileUpdateLogic {
	return &ProfileUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProfileUpdateLogic) ProfileUpdate(in *iam.ProfileUpdateRequest) (*iam.Empty, error) {
	if in == nil || in.UserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}

	userInfo, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, in.UserId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "获取用户信息失败", err))
	}

	if in.Nickname != "" {
		userInfo.Nickname = in.Nickname
	}
	if in.Avatar != "" {
		userInfo.Avatar = in.Avatar
	}
	if in.Signature != "" {
		userInfo.Signature = in.Signature
	}

	userInfo.UpdatedAt = time.Now().Unix()

	if err := l.svcCtx.Domain.IAM.User.Update(l.ctx, userInfo); err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "更新个人信息失败", err))
	}

	return &iam.Empty{}, nil
}
