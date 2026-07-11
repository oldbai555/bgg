//go:build integration

package integration

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	authlogic "postapocgame/admin-server/internal/logic/iam/auth"
	iammodel "postapocgame/admin-server/internal/model/iam"
	"postapocgame/admin-server/internal/types"
)

// 08-testing-strategy.md §5 场景 1：登录 e2e —— 真实建用户 → 真实登录 → 拿到 token。
func TestIntegration_Login_E2E(t *testing.T) {
	env := setupTestEnv(t)
	ctx := context.Background()

	username := "it_login_" + uniqueSuffix()
	password := "S3cret!Pwd"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &iammodel.AdminUser{
		Username:     username,
		Nickname:     "集成测试用户",
		PasswordHash: string(hash),
		Status:       1,
	}
	require.NoError(t, env.Domain.IAM.User.Create(ctx, user))
	require.NotZero(t, user.Id, "Create 之后应该拿到真实自增 ID")

	httpReq := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	l := authlogic.NewLoginLogic(httpReq.Context(), env.SvcCtx)
	resp, err := l.Login(&types.LoginReq{Username: username, Password: password}, httpReq)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}
