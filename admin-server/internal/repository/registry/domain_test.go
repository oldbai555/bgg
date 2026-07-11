package registry

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

	iamrepo "postapocgame/admin-server/internal/repository/iam"

	"postapocgame/admin-server/internal/repository"
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

// TestIamUserListerAdapter_FindChunk_FiltersDisabledAndDeleted 验证 registry 包内的
// iamUserListerAdapter 正确过滤掉已禁用/已删除的用户——这是"状态正常"这一 IAM 业务语义
// 留在 registry 组合根、不泄漏进 chatdomain 包的地方，值得单独测。
func TestIamUserListerAdapter_FindChunk_FiltersDisabledAndDeleted(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepo(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).
			AddRow(1, "alice", "", "", "", "", 0, int64(1), 0, 0, 0).  // 启用，保留
			AddRow(2, "bob", "", "", "", "", 0, int64(0), 0, 0, 0).    // 禁用，过滤
			AddRow(3, "carol", "", "", "", "", 0, int64(1), 0, 0, 99)) // 已删除，过滤

	adapter := &iamUserListerAdapter{userRepo: iamrepo.NewUserRepository(repo)}
	refs, lastID, err := adapter.FindChunk(context.Background(), 20, 0)

	require.NoError(t, err)
	assert.Equal(t, uint64(3), lastID) // 分页游标以底层返回的最后一条记录 ID 为准，不受客户端过滤影响
	require.Len(t, refs, 1)
	assert.Equal(t, uint64(1), refs[0].ID)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
