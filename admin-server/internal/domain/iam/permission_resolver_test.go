package iam_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	iamdomain "postapocgame/admin-server/internal/domain/iam"
)

var adminApiColumns = []string{"id", "name", "method", "path", "description", "status", "created_at", "updated_at", "deleted_at"}

func TestPermissionResolver_CanAccess_SuperAdminBypass(t *testing.T) {
	repo, _, _, cleanup := newTestRepo(t)
	defer cleanup()

	resolver := iamdomain.NewPermissionResolver(repo)
	// userID == 1 是超级管理员，不查任何表直接放行。
	allowed, err := resolver.CanAccess(context.Background(), 1, "GET", "/api/v1/iam/user/list")

	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestPermissionResolver_CanAccess_ApiNotFound(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_api`")).
		WillReturnError(sql.ErrNoRows)

	resolver := iamdomain.NewPermissionResolver(repo)
	allowed, err := resolver.CanAccess(context.Background(), 2, "GET", "/api/v1/not-exist")

	require.Error(t, err)
	assert.False(t, allowed)
}

func TestPermissionResolver_CanAccess_NoRoles_Denied(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_api`")).
		WillReturnRows(sqlmock.NewRows(adminApiColumns).
			AddRow(10, "用户列表", "GET", "/api/v1/iam/user/list", "", 1, 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("from admin_user_role")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "role_id", "created_at", "updated_at"}))

	resolver := iamdomain.NewPermissionResolver(repo)
	allowed, err := resolver.CanAccess(context.Background(), 2, "GET", "/api/v1/iam/user/list")

	require.NoError(t, err)
	assert.False(t, allowed)
}

func TestPermissionResolver_CanAccess_MatchedApiPermission_Allowed(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_api`")).
		WillReturnRows(sqlmock.NewRows(adminApiColumns).
			AddRow(10, "用户列表", "GET", "/api/v1/iam/user/list", "", 1, 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("from admin_user_role")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "role_id", "created_at", "updated_at"}).AddRow(1, 2, 5, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("FROM admin_role_permission")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "permission_id", "created_at", "updated_at"}).AddRow(1, 5, 20, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("from admin_permission_api")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "permission_id", "api_id", "created_at", "updated_at"}).AddRow(1, 20, 10, 0, 0))

	resolver := iamdomain.NewPermissionResolver(repo)
	allowed, err := resolver.CanAccess(context.Background(), 2, "GET", "/api/v1/iam/user/list")

	require.NoError(t, err)
	assert.True(t, allowed)
}

func TestPermissionResolver_CanAccess_NoMatchedPermission_Denied(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_api`")).
		WillReturnRows(sqlmock.NewRows(adminApiColumns).
			AddRow(10, "用户列表", "GET", "/api/v1/iam/user/list", "", 1, 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("from admin_user_role")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "role_id", "created_at", "updated_at"}).AddRow(1, 2, 5, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("FROM admin_role_permission")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "role_id", "permission_id", "created_at", "updated_at"}).AddRow(1, 5, 20, 0, 0))
	// 目标接口 10 关联的权限是 99，用户拥有的权限集合里只有 20，两者不相交。
	sqlMock.ExpectQuery(regexp.QuoteMeta("from admin_permission_api")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "permission_id", "api_id", "created_at", "updated_at"}).AddRow(1, 99, 10, 0, 0))

	resolver := iamdomain.NewPermissionResolver(repo)
	allowed, err := resolver.CanAccess(context.Background(), 2, "GET", "/api/v1/iam/user/list")

	require.NoError(t, err)
	assert.False(t, allowed)
}
