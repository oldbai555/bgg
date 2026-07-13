package server

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

	"postapocgame/admin-server/services/iam/internal/repository"
	"postapocgame/admin-server/services/iam/internal/repository/registry"
	pb "postapocgame/admin-server/pkg/iamcallback/pb"
)

var adminUserColumns = []string{
	"id", "username", "nickname", "password_hash", "avatar", "signature",
	"department_id", "status", "created_at", "updated_at", "deleted_at",
}

func newTestServer(t *testing.T) (*IamCallbackServer, sqlmock.Sqlmock, func()) {
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

	return NewIamCallbackServer(registry.NewDomain(repo)), sqlMock, func() {
		_ = db.Close()
		mr.Close()
	}
}

// TestFindActiveUserChunk_FiltersDisabledAndDeleted 验证过滤掉已禁用/已删除的用户——这是
// "状态正常"这一 IAM 业务语义留在这里（不泄漏进 chat-rpc）的地方，原来是
// internal/repository/registry 包内 iamUserListerAdapter.FindChunk 的职责，chat 域拆分后
// 搬到这个临时的 IamCallback server 实现里，见 server.go 包注释。
func TestFindActiveUserChunk_FiltersDisabledAndDeleted(t *testing.T) {
	s, sqlMock, cleanup := newTestServer(t)
	defer cleanup()

	sqlMock.ExpectQuery(regexp.QuoteMeta("from `admin_user`")).
		WillReturnRows(sqlmock.NewRows(adminUserColumns).
			AddRow(1, "alice", "", "", "", "", 0, int64(1), 0, 0, 0).  // 启用，保留
			AddRow(2, "bob", "", "", "", "", 0, int64(0), 0, 0, 0).    // 禁用，过滤
			AddRow(3, "carol", "", "", "", "", 0, int64(1), 0, 0, 99)) // 已删除，过滤

	resp, err := s.FindActiveUserChunk(context.Background(), &pb.FindActiveUserChunkRequest{Limit: 20, LastId: 0})

	require.NoError(t, err)
	assert.Equal(t, uint64(3), resp.NextLastId) // 分页游标以底层返回的最后一条记录 ID 为准，不受客户端过滤影响
	require.Len(t, resp.Users, 1)
	assert.Equal(t, uint64(1), resp.Users[0].Id)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
