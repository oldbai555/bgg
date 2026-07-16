package iam_test

import (
	"regexp"
	"testing"

	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	iamdomain "postapocgame/admin-server/services/iam/internal/domain/iam"
	"postapocgame/admin-server/pkg/errs"
)

var adminRoleColumns = []string{"id", "name", "code", "description", "status", "created_at", "updated_at", "deleted_at"}
var adminPermissionColumns = []string{"id", "name", "code", "description", "created_at", "updated_at", "deleted_at"}

func TestRBACService_UpdateRolePermissions_HappyPath(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_role`")).
		WillReturnRows(sqlmock.NewRows(adminRoleColumns).AddRow(1, "admin", "admin", nil, int64(1), 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("admin_permission")).
		WillReturnRows(sqlmock.NewRows(adminPermissionColumns).
			AddRow(10, "p1", "p1", nil, 0, 0, 0).
			AddRow(11, "p2", "p2", nil, 0, 0, 0))
	sqlMock.ExpectExec(regexp.QuoteMeta("DELETE FROM admin_role_permission")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO admin_role_permission")).
		WillReturnResult(sqlmock.NewResult(0, 2))
	sqlMock.ExpectCommit()

	svc := iamdomain.NewRBACService(repo)
	err := svc.UpdateRolePermissions(context.Background(), 1, []uint64{10, 11})

	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRBACService_UpdateRolePermissions_RollbackOnInsertError(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_role`")).
		WillReturnRows(sqlmock.NewRows(adminRoleColumns).AddRow(1, "admin", "admin", nil, int64(1), 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("admin_permission")).
		WillReturnRows(sqlmock.NewRows(adminPermissionColumns).
			AddRow(10, "p1", "p1", nil, 0, 0, 0))
	sqlMock.ExpectExec(regexp.QuoteMeta("DELETE FROM admin_role_permission")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO admin_role_permission")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	svc := iamdomain.NewRBACService(repo)
	err := svc.UpdateRolePermissions(context.Background(), 1, []uint64{10})

	require.Error(t, err)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeInternalError, bizErr.Code)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRBACService_UpdateRolePermissions_RollbackOnRoleNotFound(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_role`")).
		WillReturnRows(sqlmock.NewRows(adminRoleColumns)) // 角色不存在
	sqlMock.ExpectRollback()

	svc := iamdomain.NewRBACService(repo)
	err := svc.UpdateRolePermissions(context.Background(), 999, []uint64{10})

	require.Error(t, err)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeBadRequest, bizErr.Code)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
