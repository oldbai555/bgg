//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/logic"
	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
)

// 08-testing-strategy.md §5 场景 1：登录 e2e —— 真实建用户 → 真实登录 → 拿到 token。
// Login 的业务逻辑（用户查找/密码校验/生成 token/登录日志/未读公告通知）已经从 gateway
// 整段搬进 iam-rpc 的 LoginLogic，这里直接测 iam-rpc 自己的 Logic，不再经过 gateway。
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

	l := logic.NewLoginLogic(ctx, env.SvcCtx)
	resp, err := l.Login(&iam.LoginRequest{Username: username, Password: password})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}
