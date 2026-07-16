package repository_test

import (
	"context"
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

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/repository"
)

// newTestRepo 用 sqlmock 打桩 DB、miniredis 打桩 CachedConn 依赖的 Redis，
// 构造一个可以真实验证 Repository.Transact 事务语义的 *repository.Repository。
func newTestRepo(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
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

	return repo, mock, func() {
		_ = db.Close()
		mr.Close()
	}
}

func TestRepository_Transact_HappyPathCommits(t *testing.T) {
	repo, mock, cleanup := newTestRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("insert into")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Transact(context.Background(), func(ctx context.Context, txRepo *repository.Repository) error {
		// txRepo 是换绑过事务 session 的克隆，AdminUserModel.Insert 必须走同一个事务连接，
		// 不能继续闭包引用外层 repo —— 这里通过参数 txRepo 而不是 repo 调用即验证了这一点。
		_, err := txRepo.AdminUserModel.Insert(ctx, &iammodel.AdminUser{Username: "alice"})
		return err
	})

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_Transact_RollbackOnError(t *testing.T) {
	repo, mock, cleanup := newTestRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("insert into")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectRollback()

	wantErr := errors.New("boom after insert")
	err := repo.Transact(context.Background(), func(ctx context.Context, txRepo *repository.Repository) error {
		if _, err := txRepo.AdminUserModel.Insert(ctx, &iammodel.AdminUser{Username: "bob"}); err != nil {
			return err
		}
		// 模拟事务内第二步失败：Insert 本身已经成功执行到 DB，但整体方法必须回滚，
		// 验证的正是"部分执行后回滚"这个真实故障模式，而不是第一步就失败的简单情形。
		return wantErr
	})

	require.ErrorIs(t, err, wantErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}
