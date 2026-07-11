package iam_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	iamdomain "postapocgame/admin-server/internal/domain/iam"
	chatmocks "postapocgame/admin-server/internal/mocks/chat"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/pkg/errs"
)

var adminUserColumns = []string{
	"id", "username", "nickname", "password_hash", "avatar", "signature",
	"department_id", "status", "created_at", "updated_at", "deleted_at",
}

func newTestRepo(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	conn := sqlx.NewSqlConnFromDB(db)

	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisConf := redis.RedisConf{Host: mr.Addr(), Type: "node"}
	cacheConf := cache.CacheConf{{RedisConf: redisConf, Weight: 100}}
	rdb, err := redis.NewRedis(redisConf)
	require.NoError(t, err)

	repo, err := repository.NewRepository(conn, cacheConf, rdb)
	require.NoError(t, err)

	return repo, sqlMock, func() {
		_ = db.Close()
		mr.Close()
	}
}

func TestUserDomainService_CreateUser_HappyPath(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns)) // 用户名查询返回空结果集
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_user`")).
		WillReturnResult(sqlmock.NewResult(42, 1))
	sqlMock.ExpectCommit()

	onboarding := chatmocks.NewOnboarding(t)
	done := make(chan uint64, 1)
	onboarding.On("InitNewUser", mock.Anything, uint64(42)).
		Run(func(args mock.Arguments) { done <- args.Get(1).(uint64) }).
		Return(nil)

	svc := iamdomain.NewUserDomainService(repo, onboarding)
	user, err := svc.CreateUser(context.Background(), iamdomain.CreateUserInput{
		Username: "alice",
		Password: "s3cret",
	})

	require.NoError(t, err)
	assert.Equal(t, uint64(42), user.Id)
	assert.NoError(t, sqlMock.ExpectationsWereMet())

	select {
	case gotUserID := <-done:
		assert.Equal(t, uint64(42), gotUserID)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for async chat onboarding to be triggered")
	}
}

func TestUserDomainService_CreateUser_RollbackOnDuplicateUsername(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).AddRow(
			1, "alice", "", "hash", "", "", 0, int64(1), 0, 0, 0,
		)) // 用户名已存在
	sqlMock.ExpectRollback()

	onboarding := chatmocks.NewOnboarding(t) // 不应该被调用

	svc := iamdomain.NewUserDomainService(repo, onboarding)
	_, err := svc.CreateUser(context.Background(), iamdomain.CreateUserInput{
		Username: "alice",
		Password: "s3cret",
	})

	require.Error(t, err)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeBadRequest, bizErr.Code)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestUserDomainService_CreateUser_RollbackOnInsertError(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_user`")).
		WillReturnError(errors.New("db down"))
	sqlMock.ExpectRollback()

	onboarding := chatmocks.NewOnboarding(t) // 不应该被调用

	svc := iamdomain.NewUserDomainService(repo, onboarding)
	_, err := svc.CreateUser(context.Background(), iamdomain.CreateUserInput{
		Username: "bob",
		Password: "s3cret",
	})

	require.Error(t, err)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeInternalError, bizErr.Code)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
