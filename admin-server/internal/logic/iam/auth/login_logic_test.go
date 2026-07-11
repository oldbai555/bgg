package auth

import (
	"net/http/httptest"
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

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/repository/registry"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
)

var adminUserColumns = []string{
	"id", "username", "nickname", "password_hash", "avatar", "signature",
	"department_id", "status", "created_at", "updated_at", "deleted_at",
}

// newTestSvcCtx 用 sqlmock 打桩 DB、miniredis 打桩 CachedConn 依赖的 Redis，构造一个
// 可以跑 Login 的最小 svc.ServiceContext（Repository + Domain + JWT 配置）。
// recordLoginLog/createUnreadNoticeNotifications 是登录流程里既有的异步尽力而为写入，
// 不在本文件的断言范围内（未显式 mock 时命中 sqlmock 的“未预期调用”只会被内部 logx 记录，不会让测试失败）。
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
			JWT: config.JWTConf{
				AccessSecret:  "test-access-secret",
				RefreshSecret: "test-refresh-secret",
				AccessExpire:  3600,
				RefreshExpire: 86400,
				Issuer:        "admin-server-test",
			},
		},
	}
	return svcCtx, sqlMock, mr, func() {
		_ = db.Close()
		mr.Close()
	}
}

func TestLoginLogic_Login_UserNotFound(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnError(sqlmock.ErrCancelled)

	httpReq := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	l := NewLoginLogic(httpReq.Context(), svcCtx)
	resp, err := l.Login(&types.LoginReq{Username: "nobody", Password: "whatever"}, httpReq)

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

	httpReq := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	l := NewLoginLogic(httpReq.Context(), svcCtx)
	resp, err := l.Login(&types.LoginReq{Username: "alice", Password: "wrong-password"}, httpReq)

	require.Error(t, err)
	assert.Nil(t, resp)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeUnauthorized, bizErr.Code)
}

func TestLoginLogic_Login_UserDisabled(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).
			AddRow(1, "alice", "Alice", string(hash), "", "", 0, 0, 0, 0, 0)) // status=0 已禁用

	httpReq := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	l := NewLoginLogic(httpReq.Context(), svcCtx)
	resp, err := l.Login(&types.LoginReq{Username: "alice", Password: "correct-password"}, httpReq)

	require.Error(t, err)
	assert.Nil(t, resp)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeForbidden, bizErr.Code)
}

func TestLoginLogic_Login_Success(t *testing.T) {
	svcCtx, sqlMock, _, cleanup := newTestSvcCtx(t)
	defer cleanup()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).
			AddRow(1, "alice", "Alice", string(hash), "", "", 0, 1, 0, 0, 0))

	httpReq := httptest.NewRequest("POST", "/api/v1/auth/login", nil)
	l := NewLoginLogic(httpReq.Context(), svcCtx)
	resp, err := l.Login(&types.LoginReq{Username: "alice", Password: "correct-password"}, httpReq)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)

	// recordLoginLog/createUnreadNoticeNotifications 是登录成功后触发的异步尽力而为写入，
	// 给它们一点时间跑完，避免 goroutine 在 cleanup() 关闭 DB 之后才执行触发无关的日志噪音。
	time.Sleep(20 * time.Millisecond)
}
