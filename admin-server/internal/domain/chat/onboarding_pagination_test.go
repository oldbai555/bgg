package chat_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	chatdomain "postapocgame/admin-server/internal/domain/chat"
	chatmocks "postapocgame/admin-server/internal/mocks/chat"
	"postapocgame/admin-server/internal/repository"
)

var chatColumns = []string{"id", "name", "type", "avatar", "description", "created_by", "created_at", "updated_at", "deleted_at"}

func newTestRepoForPaginationTest(t *testing.T) (*repository.Repository, sqlmock.Sqlmock, func()) {
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

// TestChatOnboardingService_InitNewUser_PaginatesUntilShortPage 通过导出的 InitNewUser 入口
// 间接验证 createPrivateChatsForExistingUsers 的分页游标：page1 恰好 limit(100) 条触发第二次
// FindChunk 调用，page2 不足 limit 提前结束循环。用 mockery 生成的 UserLister mock 断言调用
// 次数和 lastID 参数，不需要关心 IAM 那边的 SQL 长什么样——这是窄接口带来的直接好处。
func TestChatOnboardingService_InitNewUser_PaginatesUntilShortPage(t *testing.T) {
	repo, sqlMock, cleanup := newTestRepoForPaginationTest(t)
	defer cleanup()

	// 1. joinDefaultGroup：默认企业群组查询返回空结果集，提前返回，不再产生后续 SQL。
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `chat`")).
		WillReturnRows(sqlmock.NewRows(chatColumns))

	// 2. createPrivateChatsForExistingUsers：两页存量用户。
	const pageSize = 100
	page1 := make([]chatdomain.UserRef, 0, pageSize)
	for i := uint64(0); i < pageSize; i++ {
		page1 = append(page1, chatdomain.UserRef{ID: 100 + i})
	}
	page2 := []chatdomain.UserRef{{ID: 300}}

	lister := chatmocks.NewUserLister(t)
	lister.On("FindChunk", mock.Anything, 100, uint64(0)).Return(page1, uint64(199), nil).Once()
	lister.On("FindChunk", mock.Anything, 100, uint64(199)).Return(page2, uint64(0), nil).Once()

	// 每个存量用户都命中"私聊已存在"分支，避免再展开一整套建私聊事务的 SQL 期望——
	// 本测试关心的是分页游标是否正确传递，不是 createPrivateChat 本身（后者见 onboarding_test.go）。
	existingChatRow := func() *sqlmock.Rows {
		return sqlmock.NewRows(chatColumns).AddRow(1, "", 1, "", "", 0, 0, 0, 0)
	}
	for i := 0; i < len(page1)+len(page2); i++ {
		sqlMock.ExpectQuery(regexp.QuoteMeta("FROM chat c")).WillReturnRows(existingChatRow())
	}

	svc := chatdomain.NewChatOnboardingService(repo, lister)
	err := svc.InitNewUser(context.Background(), 1)

	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
	lister.AssertExpectations(t)
}
