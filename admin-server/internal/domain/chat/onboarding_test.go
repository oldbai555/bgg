package chat

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/internal/repository"
)

// newTestRepo 用 sqlmock 打桩 DB、miniredis 打桩 CachedConn 依赖的 Redis。
// 白盒测试（package chat），不能 import internal/mocks/chat——该包反过来 import 本包，
// 会形成 test 场景下的 import cycle，这也是本文件只测不依赖 UserLister 的 createPrivateChat 的原因；
// 需要 UserLister mock 的分页测试见 onboarding_pagination_test.go（外部测试包 chat_test）。
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

func TestChatOnboardingService_CreatePrivateChat_HappyPath(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
		WillReturnResult(sqlmock.NewResult(7, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(2, 1))
	sqlMock.ExpectCommit()

	svc := NewChatOnboardingService(repo, nil) // createPrivateChat 不使用 userLister
	err := svc.createPrivateChat(context.Background(), 1, 2)

	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestChatOnboardingService_CreatePrivateChat_RollbackOnSecondChatUserInsertError(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
		WillReturnResult(sqlmock.NewResult(7, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	svc := NewChatOnboardingService(repo, nil)
	err := svc.createPrivateChat(context.Background(), 1, 2)

	require.Error(t, err)
	// sqlmock.ExpectationsWereMet 本身即验证第三个 exec（第二条 chat_user）之后没有多余/遗漏的调用。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
