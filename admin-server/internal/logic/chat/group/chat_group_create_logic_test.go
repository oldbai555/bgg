package group

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
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	jwthelper "postapocgame/admin-server/pkg/jwt"
)

// newTestSvcCtx 用 sqlmock 打桩 DB、miniredis 打桩 CachedConn 依赖的 Redis，
// 只填充 ChatGroupCreate 实际用到的 svc.ServiceContext.Repository 字段。
func newTestSvcCtx(t *testing.T) (*svc.ServiceContext, sqlmock.Sqlmock, func()) {
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

	return &svc.ServiceContext{Repository: repo}, sqlMock, func() {
		_ = db.Close()
		mr.Close()
	}
}

func TestChatGroupCreateLogic_ChatGroupCreate_HappyPath(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
		WillReturnResult(sqlmock.NewResult(7, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()

	ctx := jwthelper.WithAuthUser(context.Background(), jwthelper.AuthUser{UserID: 1, Username: "alice"})
	l := NewChatGroupCreateLogic(ctx, svcCtx)
	resp, err := l.ChatGroupCreate(&types.ChatGroupCreateReq{Name: "测试群组"})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestChatGroupCreateLogic_ChatGroupCreate_RollbackWhenCreatorJoinFails(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
		WillReturnResult(sqlmock.NewResult(7, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	ctx := jwthelper.WithAuthUser(context.Background(), jwthelper.AuthUser{UserID: 1, Username: "alice"})
	l := NewChatGroupCreateLogic(ctx, svcCtx)
	resp, err := l.ChatGroupCreate(&types.ChatGroupCreateReq{Name: "测试群组"})

	require.Error(t, err)
	assert.Nil(t, resp)
	// 断言：群组那条 INSERT 虽然执行成功过，但整个事务未提交，不会留下没有任何成员的孤儿群组。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
