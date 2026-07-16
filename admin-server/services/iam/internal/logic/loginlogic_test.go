package logic

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/config"
	"postapocgame/admin-server/services/iam/internal/repository"
	"postapocgame/admin-server/services/iam/internal/repository/registry"
	"postapocgame/admin-server/services/iam/internal/svc"
)

var adminUserColumns = []string{
	"id", "username", "nickname", "password_hash", "avatar", "signature",
	"department_id", "status", "created_at", "updated_at", "deleted_at",
}

// newTestSvcCtx 用 sqlmock 打桩 DB、miniredis 打桩 CachedConn 依赖的 Redis，构造一个
// 可以跑 Login/Refresh 的最小 svc.ServiceContext（Repository + Domain + JWT 配置）。
// recordLoginLog/createUnreadNoticeNotifications 是登录流程里既有的异步尽力而为写入，
// 不在本文件的断言范围内。
func newTestSvcCtx(t *testing.T) (*svc.ServiceContext, sqlmock.Sqlmock, *miniredis.Miniredis, func()) {
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

	svcCtx := &svc.ServiceContext{
		Repository: repo,
		Domain:     registry.NewDomain(repo),
		Config: config.Config{
			JWT: struct {
				AccessSecret  string
				RefreshSecret string
				AccessExpire  int64
				RefreshExpire int64
				Issuer        string
			}{
				AccessSecret:  "test-access-secret",
				RefreshSecret: "test-refresh-secret",
				AccessExpire:  3600,
				RefreshExpire: 86400,
				Issuer:        "iam-rpc-test",
			},
		},
	}
	return svcCtx, sqlMock, mr, func() {
		_ = db.Close()
		mr.Close()
	}
}

func requireGRPCCode(t *testing.T, err error, code codes.Code) {
	t.Helper()
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "error should be a grpc status error: %v", err)
	assert.Equal(t, code, st.Code())
}

func TestLoginLogic_Login_UserNotFound(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnError(sqlmock.ErrCancelled)

	l := NewLoginLogic(context.Background(), svcCtx)
	resp, err := l.Login(&iam.LoginRequest{Username: "nobody", Password: "whatever"})

	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestLoginLogic_Login_WrongPassword(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).
			AddRow(1, "alice", "Alice", string(hash), "", "", 0, 1, 0, 0, 0))

	l := NewLoginLogic(context.Background(), svcCtx)
	resp, err := l.Login(&iam.LoginRequest{Username: "alice", Password: "wrong-password"})

	assert.Nil(t, resp)
	requireGRPCCode(t, err, codes.Unauthenticated)
}

func TestLoginLogic_Login_UserDisabled(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).
			AddRow(1, "alice", "Alice", string(hash), "", "", 0, 0, 0, 0, 0)) // status=0 已禁用

	l := NewLoginLogic(context.Background(), svcCtx)
	resp, err := l.Login(&iam.LoginRequest{Username: "alice", Password: "correct-password"})

	assert.Nil(t, resp)
	requireGRPCCode(t, err, codes.PermissionDenied)
}

func TestLoginLogic_Login_Success(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).
			AddRow(1, "alice", "Alice", string(hash), "", "", 0, 1, 0, 0, 0))

	l := NewLoginLogic(context.Background(), svcCtx)
	resp, err := l.Login(&iam.LoginRequest{Username: "alice", Password: "correct-password"})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	// recordLoginLog/createUnreadNoticeNotifications 是登录成功后触发的异步尽力而为写入，
	// 给它们一点时间跑完，避免 goroutine 在 cleanup() 关闭 DB 之后才执行触发无关的日志噪音。
	time.Sleep(20 * time.Millisecond)
}
