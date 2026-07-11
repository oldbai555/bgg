//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"postapocgame/admin-server/internal/middleware"
	iammodel "postapocgame/admin-server/internal/model/iam"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
)

// buildProtectedHandler 组一条 AuthMiddleware -> PermissionMiddleware -> 200 的最小请求链，
// 和 admin.api 里 middleware: Auth,Permission 的声明顺序一致（见 00-workflow.md 中间件顺序要求）。
// 返回的 *bool 记录下游 handler 是否被放行调用到，供拒绝场景断言"没有放行"。
func buildProtectedHandler(t *testing.T, env *testEnv) (http.HandlerFunc, *bool) {
	t.Helper()
	authMw := middleware.NewAuthMiddleware(env.SvcCtx.Config, env.Repo)
	permMw := middleware.NewPermissionMiddleware(env.Domain.IAM.PermissionResolver)

	called := false
	final := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}
	return authMw.Handle(permMw.Handle(final)), &called
}

// createTestUserWithToken 建一个真实用户并签发一个真实 access token，返回用户和 token。
func createTestUserWithToken(t *testing.T, env *testEnv, usernamePrefix string) (*iammodel.AdminUser, string) {
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

	token, err := jwthelper.GenerateToken(
		env.SvcCtx.Config.JWT.AccessSecret, env.SvcCtx.Config.JWT.Issuer,
		env.SvcCtx.Config.JWT.AccessExpire, user.Id, user.Username, false)
	require.NoError(t, err)

	return user, token
}

// 08-testing-strategy.md §5 场景 2：RBAC 允许的请求——有权限的用户访问受保护接口，返回 200。
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

	user, token := createTestUserWithToken(t, env, "it_rbac_allow")
	require.NoError(t, env.Domain.IAM.UserRole.UpdateUserRoles(ctx, user.Id, []uint64{role.Id}))

	handler, called := buildProtectedHandler(t, env)
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, *called, "有权限时应该放行到下一个 handler")
}

// 08-testing-strategy.md §5 场景 3：RBAC 拒绝的请求——无权限用户访问受保护接口，返回 403（业务码）。
func TestIntegration_RBAC_Denied(t *testing.T) {
	env := setupTestEnv(t)
	ctx := context.Background()

	const method = "GET"
	path := "/api/v1/integration-test/rbac-denied-" + uniqueSuffix()

	api := &iammodel.AdminApi{Name: "集成测试-拒绝", Method: method, Path: path, Status: 1}
	require.NoError(t, env.Domain.IAM.Api.Create(ctx, api))

	// 用户没有任何角色，CanAccess 在 ListRoleIDsByUserID 返回空之后直接判定 false。
	_, token := createTestUserWithToken(t, env, "it_rbac_deny")

	handler, called := buildProtectedHandler(t, env)
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler(rec, req)

	// pkg/response.ErrorCtx 对业务错误统一写 HTTP 400，业务错误码在响应体里，
	// 与 internal/middleware/permissionmiddleware_test.go 的断言口径一致。
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, *called, "无权限时不应该放行到下一个 handler")

	var body struct {
		Code int `json:"code"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, errs.CodeForbidden, body.Code)
}
