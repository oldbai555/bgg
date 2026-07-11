//go:build integration

// Package integration 是 08-testing-strategy.md §5 要求的集成测试套件：跑真实 MySQL + Redis，
// 不随 `go test ./...` 默认触发，只有显式 `go test -tags=integration ./...` 才会编译进这个包。
// 数据库/Redis 地址通过环境变量传入（CI 见 .github/workflows/ci.yml 的 integration-test job），
// 本地没有配置这两个环境变量时用 t.Skip 跳过，而不是报错，避免开发者本地没有 MySQL/Redis 时跑不过。
package integration

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/repository/registry"
	"postapocgame/admin-server/internal/svc"
)

// testEnv 聚合一次集成测试需要的真实基础设施连接。
type testEnv struct {
	Repo   *repository.Repository
	Domain *registry.Domain
	SvcCtx *svc.ServiceContext
}

// setupTestEnv 用 TEST_MYSQL_DSN/TEST_REDIS_ADDR 连接真实 MySQL + Redis。
// 两个环境变量任一缺失就 t.Skip，不影响没有本地 MySQL/Redis 的开发者跑其余的单元测试。
func setupTestEnv(t *testing.T) *testEnv {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	if dsn == "" || redisAddr == "" {
		t.Skip("TEST_MYSQL_DSN/TEST_REDIS_ADDR 未设置，跳过集成测试（本地需要先起 docker-compose 或连接开发库）")
	}

	conn := sqlx.NewMysql(dsn)
	// 集成测试库 admin_test 每次运行前由 CI/docker-compose 按 db/docker-init.sh 的顺序重新初始化，
	// 探测一次连接，连不上就快速失败而不是让后续每条 SQL 都超时。
	rawDB, err := conn.RawDB()
	require.NoError(t, err)
	require.NoError(t, rawDB.Ping())

	rdb, err := redis.NewRedis(redis.RedisConf{Host: redisAddr, Type: "node"})
	require.NoError(t, err)

	cacheConf := cache.CacheConf{{RedisConf: redis.RedisConf{Host: redisAddr, Type: "node"}, Weight: 100}}
	repo, err := repository.NewRepository(conn, cacheConf, rdb)
	require.NoError(t, err)

	domain := registry.NewDomain(repo)

	svcCtx := &svc.ServiceContext{
		Repository: repo,
		Domain:     domain,
		Config: config.Config{
			JWT: config.JWTConf{
				AccessSecret:  "integration-test-access-secret",
				RefreshSecret: "integration-test-refresh-secret",
				AccessExpire:  3600,
				RefreshExpire: 86400,
				Issuer:        "admin-server-integration-test",
			},
		},
	}

	return &testEnv{Repo: repo, Domain: domain, SvcCtx: svcCtx}
}

// uniqueSuffix 用当前纳秒时间戳拼出一个测试内唯一后缀，避免重复运行时 username/code 撞唯一索引。
func uniqueSuffix() string {
	return time.Now().Format("150405.000000000")
}
