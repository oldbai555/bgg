package sdk_test

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

	sdkdomain "postapocgame/admin-server/services/sdk/internal/domain/sdk"
	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
	"postapocgame/admin-server/services/sdk/internal/repository"
)

// newTestStore 和 internal/repository/repository_test.go 用的是同一套组合：miniredis +
// 单节点 CacheConf——goctl 生成的 Model 内部走 sqlc.CachedConn，cache.New 对空 CacheConf
// 会 log.Fatal，必须喂一个真实（哪怕是内存模拟的）Redis 节点。
func newTestStore(t *testing.T) (*repository.Store, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	conn := sqlx.NewSqlConnFromDB(db)

	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisConf := redis.RedisConf{Host: mr.Addr(), Type: "node"}
	cacheConf := cache.CacheConf{{RedisConf: redisConf, Weight: 100}}

	store := repository.NewStore(conn, cacheConf)

	return store, sqlMock, func() {
		_ = db.Close()
		mr.Close()
	}
}

func TestSDKService_SaveApiKeyBindings_HappyPath(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE sdk_key_api")).
		WillReturnResult(sqlmock.NewResult(0, 2))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `sdk_key_api`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `sdk_key_api`")).
		WillReturnResult(sqlmock.NewResult(2, 1))
	sqlMock.ExpectCommit()

	svc := sdkdomain.NewSDKService(store)
	bindings := []sdkmodel.SdkKeyApi{
		{SdkKeyId: 1, SdkInterfaceId: 10},
		{SdkKeyId: 1, SdkInterfaceId: 11},
	}
	err := svc.SaveApiKeyBindings(context.Background(), 1, bindings)

	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestSDKService_SaveApiKeyBindings_RollbackOnInsertError(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE sdk_key_api")).
		WillReturnResult(sqlmock.NewResult(0, 2))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `sdk_key_api`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `sdk_key_api`")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	svc := sdkdomain.NewSDKService(store)
	bindings := []sdkmodel.SdkKeyApi{
		{SdkKeyId: 1, SdkInterfaceId: 10},
		{SdkKeyId: 1, SdkInterfaceId: 11},
	}
	err := svc.SaveApiKeyBindings(context.Background(), 1, bindings)

	require.Error(t, err)
	// DELETE 和第一条 INSERT 虽然执行成功过，但整个事务未提交，随 Rollback 一起撤销，
	// 旧绑定"看起来"没有被清空。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
