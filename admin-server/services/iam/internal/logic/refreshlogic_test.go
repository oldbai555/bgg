package logic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	"postapocgame/admin-server/services/iam/iam"
	iamconsts "postapocgame/admin-server/services/iam/internal/consts"

	jwthelper "postapocgame/admin-server/pkg/jwt"
)

func TestRefreshLogic_Refresh_TokenExpired(t *testing.T) {
	svcCtx, _, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	// expireSeconds = -10：生成时即已过期的刷新令牌。
	expiredToken, err := jwthelper.GenerateToken(
		svcCtx.Config.JWT.RefreshSecret, svcCtx.Config.JWT.Issuer, -10, 1, "alice", true)
	require.NoError(t, err)

	l := NewRefreshLogic(context.Background(), svcCtx)
	resp, err := l.Refresh(&iam.RefreshRequest{RefreshToken: expiredToken})

	assert.Nil(t, resp)
	requireGRPCCode(t, err, codes.Unauthenticated)
}

func TestRefreshLogic_Refresh_TokenBlacklisted(t *testing.T) {
	svcCtx, _, mr, cleanup := newTestSvcCtx(t)
	defer cleanup()

	refreshToken, err := jwthelper.GenerateToken(
		svcCtx.Config.JWT.RefreshSecret, svcCtx.Config.JWT.Issuer, svcCtx.Config.JWT.RefreshExpire, 1, "alice", true)
	require.NoError(t, err)

	// 模拟该刷新令牌此前已经因为 Logout 被加入黑名单（直接在 miniredis 里种入黑名单 key）。
	require.NoError(t, mr.Set(iamconsts.RedisJWTBlacklistPrefix+refreshToken, "1"))

	l := NewRefreshLogic(context.Background(), svcCtx)
	resp, err := l.Refresh(&iam.RefreshRequest{RefreshToken: refreshToken})

	assert.Nil(t, resp)
	requireGRPCCode(t, err, codes.Unauthenticated)
}

func TestRefreshLogic_Refresh_Success(t *testing.T) {
	svcCtx, _, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	refreshToken, err := jwthelper.GenerateToken(
		svcCtx.Config.JWT.RefreshSecret, svcCtx.Config.JWT.Issuer, svcCtx.Config.JWT.RefreshExpire, 1, "alice", true)
	require.NoError(t, err)

	l := NewRefreshLogic(context.Background(), svcCtx)
	resp, err := l.Refresh(&iam.RefreshRequest{RefreshToken: refreshToken})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}
