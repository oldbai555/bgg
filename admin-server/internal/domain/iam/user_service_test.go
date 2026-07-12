package iam_test

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	iamdomain "postapocgame/admin-server/internal/domain/iam"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/pkg/errs"
)

var adminUserColumns = []string{
	"id", "username", "nickname", "password_hash", "avatar", "signature",
	"department_id", "status", "created_at", "updated_at", "deleted_at",
}

func newTestRepo(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, *miniredis.Miniredis, func()) {
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

	return repo, sqlMock, mr, func() {
		_ = db.Close()
		mr.Close()
	}
}

// TestUserDomainService_CreateUser_HappyPath 验证建用户成功后同步发布了
// stream:chat.user.created 事件（chat 域拆分成 chat-rpc 后的新触发机制，见
// internal/domain/iam/user_service.go 的 publishChatUserCreated）——不再需要像拆分前那样
// 轮询等待异步 goroutine 完成，XAdd 本身是同步调用。
func TestUserDomainService_CreateUser_HappyPath(t *testing.T) {
	repo, sqlMock, mr, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns)) // 用户名查询返回空结果集
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_user`")).
		WillReturnResult(sqlmock.NewResult(42, 1))
	sqlMock.ExpectCommit()

	svc := iamdomain.NewUserDomainService(repo, repo.Redis)
	user, err := svc.CreateUser(context.Background(), iamdomain.CreateUserInput{
		Username: "alice",
		Password: "s3cret",
	})

	require.NoError(t, err)
	assert.Equal(t, uint64(42), user.Id)
	assert.NoError(t, sqlMock.ExpectationsWereMet())

	entries, err := mr.Stream(iamdomain.StreamChatUserCreated)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	var event struct {
		UserID uint64 `json:"userId"`
	}
	require.NoError(t, json.Unmarshal([]byte(entries[0].Values[1]), &event)) // Values: ["payload", <json>]
	assert.Equal(t, uint64(42), event.UserID)
}

func TestUserDomainService_CreateUser_RollbackOnDuplicateUsername(t *testing.T) {
	repo, sqlMock, mr, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).AddRow(
			1, "alice", "", "hash", "", "", 0, int64(1), 0, 0, 0,
		)) // 用户名已存在
	sqlMock.ExpectRollback()

	svc := iamdomain.NewUserDomainService(repo, repo.Redis)
	_, err := svc.CreateUser(context.Background(), iamdomain.CreateUserInput{
		Username: "alice",
		Password: "s3cret",
	})

	require.Error(t, err)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeBadRequest, bizErr.Code)
	assert.NoError(t, sqlMock.ExpectationsWereMet())

	entries, err := mr.Stream(iamdomain.StreamChatUserCreated)
	require.NoError(t, err)
	assert.Empty(t, entries, "建用户失败不应该发布 chat.user.created 事件")
}

func TestUserDomainService_CreateUser_RollbackOnInsertError(t *testing.T) {
	repo, sqlMock, _, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `admin_user`")).
		WillReturnError(errors.New("db down"))
	sqlMock.ExpectRollback()

	svc := iamdomain.NewUserDomainService(repo, repo.Redis)
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
