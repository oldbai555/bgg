package logic_test

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

	"postapocgame/admin-server/services/sdk/internal/config"
	sdkdomain "postapocgame/admin-server/services/sdk/internal/domain/sdk"
	"postapocgame/admin-server/services/sdk/internal/logic"
	"postapocgame/admin-server/services/sdk/internal/repository"
	sdkrepo "postapocgame/admin-server/services/sdk/internal/repository/sdk"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"
)

// newTestSvcCtx 和 domain 包的 newTestStore 用的是同一套组合（miniredis + 单节点
// CacheConf，cache.New 对空 CacheConf 会 log.Fatal），直接拼出 VerifyApiKeyLogic 等
// 需要的完整 *svc.ServiceContext，不单独导出 store 构造函数（logic 包只应通过
// ServiceContext 访问，和其余 logic 文件的调用方式保持一致）。
func newTestSvcCtx(t *testing.T, rateLimitDefault int64) (*svc.ServiceContext, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	conn := sqlx.NewSqlConnFromDB(db)

	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisConf := redis.RedisConf{Host: mr.Addr(), Type: "node"}
	cacheConf := cache.CacheConf{{RedisConf: redisConf, Weight: 100}}

	store := repository.NewStore(conn, cacheConf)

	svcCtx := &svc.ServiceContext{
		Config:           config.Config{RateLimitDefault: rateLimitDefault},
		Admin:            sdkrepo.NewSdkAdminRepository(store),
		Public:           sdkrepo.NewSdkRepository(store),
		Service:          sdkdomain.NewSDKService(store),
		RateLimitDefault: rateLimitDefault,
	}

	return svcCtx, sqlMock, func() {
		_ = db.Close()
		mr.Close()
	}
}

var sdkKeyColumns = []string{"id", "name", "api_key", "api_secret", "status", "expire_at", "ip_whitelist", "remark", "created_at", "updated_at", "deleted_at"}

var sdkInterfaceColumns = []string{"id", "name", "api_code", "path", "method", "rate_limit_default", "status", "remark", "created_at", "updated_at", "deleted_at"}

func TestVerifyApiKey_MissingCredentials(t *testing.T) {
	svcCtx, _, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "缺少 API Key 或 Secret", resp.Message)
}

func TestVerifyApiKey_InvalidApiKey(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnError(sqlx.ErrNotFound)

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{
		ApiKey: "unknown", ApiSecret: "secret",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "无效的 API Key", resp.Message)
}

func TestVerifyApiKey_SecretMismatch(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkKeyColumns).
		AddRow(1, "key1", "abc", "real-secret", 1, 0, "", "", 0, 0, 0))

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{
		ApiKey: "abc", ApiSecret: "wrong-secret",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "Secret 不匹配", resp.Message)
}

func TestVerifyApiKey_Disabled(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkKeyColumns).
		AddRow(1, "key1", "abc", "secret", 2, 0, "", "", 0, 0, 0))

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{
		ApiKey: "abc", ApiSecret: "secret",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "API Key 已被禁用", resp.Message)
}

func TestVerifyApiKey_Expired(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkKeyColumns).
		AddRow(1, "key1", "abc", "secret", 1, time.Now().Add(-time.Hour).Unix(), "", "", 0, 0, 0))

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{
		ApiKey: "abc", ApiSecret: "secret",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "API Key 已过期", resp.Message)
}

func TestVerifyApiKey_IPNotAllowed(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkKeyColumns).
		AddRow(1, "key1", "abc", "secret", 1, 0, "10.0.0.1", "", 0, 0, 0))

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{
		ApiKey: "abc", ApiSecret: "secret", ClientIp: "1.2.3.4",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "IP 不在白名单", resp.Message)
}

func TestVerifyApiKey_InterfaceNotOpen(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkKeyColumns).
		AddRow(1, "key1", "abc", "secret", 1, 0, "", "", 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnError(sqlx.ErrNotFound)

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{
		ApiKey: "abc", ApiSecret: "secret", Method: "GET", Path: "/video/list",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "接口未开通或已禁用", resp.Message)
}

func TestVerifyApiKey_Success(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkKeyColumns).
		AddRow(1, "key1", "abc", "secret", 1, 0, "", "", 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkInterfaceColumns).
		AddRow(10, "视频列表", "get:/video/list", "/video/list", "GET", 0, 1, "", 0, 0, 0))

	resp, err := logic.NewVerifyApiKeyLogic(context.Background(), svcCtx).VerifyApiKey(&sdk.VerifyApiKeyRequest{
		ApiKey: "abc", ApiSecret: "secret", Method: "GET", Path: "/video/list",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, uint64(1), resp.SdkKeyId)
	assert.Equal(t, uint64(10), resp.SdkInterfaceId)
	assert.Equal(t, "get:/video/list", resp.ApiCode)
}

func TestGetEffectiveRateLimit_InterfaceNotFound(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnError(sqlx.ErrNotFound)

	_, err := logic.NewGetEffectiveRateLimitLogic(context.Background(), svcCtx).GetEffectiveRateLimit(&sdk.GetEffectiveRateLimitRequest{
		SdkKeyId: 1, SdkInterfaceId: 10,
	})
	require.Error(t, err)
}

func TestGetEffectiveRateLimit_StaticDefaultFallback(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 42)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkInterfaceColumns).
		AddRow(10, "视频列表", "get:/video/list", "/video/list", "GET", 0, 1, "", 0, 0, 0))
	// FindKeyApiBinding：未绑定，返回错误（FindOneBySdkKeyIdSdkInterfaceId 命中 not found）。
	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnError(sqlx.ErrNotFound)

	resp, err := logic.NewGetEffectiveRateLimitLogic(context.Background(), svcCtx).GetEffectiveRateLimit(&sdk.GetEffectiveRateLimitRequest{
		SdkKeyId: 1, SdkInterfaceId: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(42), resp.Limit)
}

func TestGetEffectiveRateLimit_CustomBindingOverrides(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 42)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(sqlmock.NewRows(sdkInterfaceColumns).
		AddRow(10, "视频列表", "get:/video/list", "/video/list", "GET", 0, 1, "", 0, 0, 0))
	sqlMock.ExpectQuery(regexp.QuoteMeta("select")).WillReturnRows(
		sqlmock.NewRows([]string{"id", "sdk_key_id", "sdk_interface_id", "custom_rate_limit", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 1, 10, 200, 0, 0, 0))

	resp, err := logic.NewGetEffectiveRateLimitLogic(context.Background(), svcCtx).GetEffectiveRateLimit(&sdk.GetEffectiveRateLimitRequest{
		SdkKeyId: 1, SdkInterfaceId: 10,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(200), resp.Limit)
}

func TestRecordCallLog_HappyPath(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, 60)
	defer cleanup()

	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `sdk_call_log`")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := logic.NewRecordCallLogLogic(context.Background(), svcCtx).RecordCallLog(&sdk.RecordCallLogRequest{
		SdkKeyId: 1, SdkInterfaceId: 10, ApiCode: "get:/video/list", Path: "/video/list", Method: "GET", RespCode: 200,
	})
	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
