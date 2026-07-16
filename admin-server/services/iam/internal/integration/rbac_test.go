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

// createTestUser 建一个真实用户，返回用户本身。RBAC 校验现在直接测 iam-rpc 的
// CheckPermissionLogic（PermissionMiddleware 现在只是调 IamRPC.CheckPermission 的薄胶水，
// 不再需要签发 JWT / 走 HTTP 中间件链来验证权限判定本身）。
func createTestUser(t *testing.T, env *testEnv, usernamePrefix string) *iammodel.AdminUser {
	t.Helper()
	ctx := context.Background()

	hash, err := bcrypt.GenerateFromPassword([]byte("S3cret!Pwd"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &iammodel.AdminUser{
		Username:     usernamePrefix + "_" + uniqueSuffix(),
		Nickname:     "RBAC 集成测试用户",
		PasswordHash: string(hash),
		Status:       1,
	}
	require.NoError(t, env.Domain.IAM.User.Create(ctx, user))
	require.NotZero(t, user.Id)

	return user
}

// 08-testing-strategy.md §5 场景 2：RBAC 允许的请求——有权限的用户访问受保护接口，Allowed=true。
func TestIntegration_RBAC_Allowed(t *testing.T) {
	env := setupTestEnv(t)
	ctx := context.Background()

	// path 带 uniqueSuffix：admin_api 的 (method,path) 有唯一索引，测试库不是每次都从空库跑
	// （例如直接对着开发库跑集成测试），固定路径重跑会撞 Duplicate entry。
	const method = "GET"
	path := "/api/v1/integration-test/rbac-allowed-" + uniqueSuffix()

	api := &iammodel.AdminApi{Name: "集成测试-允许", Method: method, Path: path, Status: 1}
	require.NoError(t, env.Domain.IAM.Api.Create(ctx, api))

	role := &iammodel.AdminRole{Name: "集成测试角色-允许" + uniqueSuffix(), Code: "it_role_allow_" + uniqueSuffix(), Status: 1}
	require.NoError(t, env.Domain.IAM.Role.Create(ctx, role))
	// 权限 ID=1 是 PermissionResolver.CanAccess 里约定的"超级权限"，拥有即放行，
	// 不需要额外维护 permission_api 关联，和领域层 permission_resolver_test.go 的
	// TestPermissionResolver_CanAccess 用例覆盖的分支保持一致的最小路径。
	require.NoError(t, env.Domain.IAM.RolePermission.UpdateRolePermissions(ctx, role.Id, []uint64{1}))

	user := createTestUser(t, env, "it_rbac_allow")
	require.NoError(t, env.Domain.IAM.UserRole.UpdateUserRoles(ctx, user.Id, []uint64{role.Id}))

	l := logic.NewCheckPermissionLogic(ctx, env.SvcCtx)
	resp, err := l.CheckPermission(&iam.CheckPermissionRequest{UserId: user.Id, Method: method, Path: path})

	require.NoError(t, err)
	assert.True(t, resp.Allowed, "有权限时 CheckPermission 应该返回 Allowed=true")
}

// 08-testing-strategy.md §5 场景 3：RBAC 拒绝的请求——无权限用户访问受保护接口，Allowed=false。
func TestIntegration_RBAC_Denied(t *testing.T) {
	env := setupTestEnv(t)
	ctx := context.Background()

	const method = "GET"
	path := "/api/v1/integration-test/rbac-denied-" + uniqueSuffix()

	api := &iammodel.AdminApi{Name: "集成测试-拒绝", Method: method, Path: path, Status: 1}
	require.NoError(t, env.Domain.IAM.Api.Create(ctx, api))

	// 用户没有任何角色，CanAccess 在 ListRoleIDsByUserID 返回空之后直接判定 false。
	user := createTestUser(t, env, "it_rbac_deny")

	l := logic.NewCheckPermissionLogic(ctx, env.SvcCtx)
	resp, err := l.CheckPermission(&iam.CheckPermissionRequest{UserId: user.Id, Method: method, Path: path})

	require.NoError(t, err)
	assert.False(t, resp.Allowed, "无权限时 CheckPermission 应该返回 Allowed=false")
}
