package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshLogic {
	return &RefreshLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshLogic) Refresh(in *iam.RefreshRequest) (*iam.TokenPair, error) {
	if in == nil || in.RefreshToken == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "刷新令牌不能为空"))
	}

	claims, err := jwthelper.ParseToken(in.RefreshToken, l.svcCtx.Config.JWT.RefreshSecret)
	if err != nil || !claims.IsRefresh {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "刷新令牌无效或已过期"))
	}

	blacklisted, err := l.svcCtx.Domain.IAM.TokenBlacklist.IsBlacklisted(l.ctx, in.RefreshToken)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "检查令牌黑名单失败", err))
	}
	if blacklisted {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "刷新令牌无效或已过期"))
	}

	accessToken, err := jwthelper.GenerateToken(
		l.svcCtx.Config.JWT.AccessSecret,
		l.svcCtx.Config.JWT.Issuer,
		l.svcCtx.Config.JWT.AccessExpire,
		claims.UserID,
		claims.Username,
		false,
	)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "生成访问令牌失败", err))
	}

	refreshToken, err := jwthelper.GenerateToken(
		l.svcCtx.Config.JWT.RefreshSecret,
		l.svcCtx.Config.JWT.Issuer,
		l.svcCtx.Config.JWT.RefreshExpire,
		claims.UserID,
		claims.Username,
		true,
	)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "生成刷新令牌失败", err))
	}

	return &iam.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
