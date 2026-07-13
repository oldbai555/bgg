package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *iam.LogoutRequest) (*iam.Empty, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	if in.AccessToken != "" {
		if err := l.svcCtx.Domain.IAM.TokenBlacklist.Blacklist(l.ctx, in.AccessToken, time.Duration(l.svcCtx.Config.JWT.AccessExpire)*time.Second); err != nil {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "加入访问令牌黑名单失败", err))
		}
	}

	if in.RefreshToken != "" {
		if err := l.svcCtx.Domain.IAM.TokenBlacklist.Blacklist(l.ctx, in.RefreshToken, time.Duration(l.svcCtx.Config.JWT.RefreshExpire)*time.Second); err != nil {
			return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "加入刷新令牌黑名单失败", err))
		}
	}

	return &iam.Empty{}, nil
}
