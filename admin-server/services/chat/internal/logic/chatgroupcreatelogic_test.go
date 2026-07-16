package logic

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

	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/repository"
	chatrepo "postapocgame/admin-server/services/chat/internal/repository/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

// newTestSvcCtx 和 services/sdk/internal/domain/sdk/sdk_service_test.go 用的是同一套组合：
// miniredis + sqlmock，只填充 ChatGroupCreate 实际用到的字段（不需要 IamCallback，本测试
// 不带初始成员）。
func newTestSvcCtx(t *testing.T) (*svc.ServiceContext, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	conn := sqlx.NewSqlConnFromDB(db)

	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisConf := redis.RedisConf{Host: mr.Addr(), Type: "node"}
	cacheConf := cache.CacheConf{{RedisConf: redisConf, Weight: 100}}
	store := repository.NewStore(conn, cacheConf)

	return &svc.ServiceContext{
			Store:    store,
			Chat:     chatrepo.NewChatRepository(store),
			ChatUser: chatrepo.NewChatUserRepository(store),
		}, sqlMock, func() {
			_ = db.Close()
			mr.Close()
		}
}

// TestChatGroupCreateLogic_ChatGroupCreate_HappyPath / _RollbackWhenCreatorJoinFails 从
// internal/logic/chat/group/chat_group_create_logic_test.go 原样迁移（chat 域拆分后网关侧
// 只剩薄胶水，Transact 相关的真正业务逻辑测试搬到这里，见 chatgroupcreatelogic.go）。
func TestChatGroupCreateLogic_ChatGroupCreate_HappyPath(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
		WillReturnResult(sqlmock.NewResult(7, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()

	l := NewChatGroupCreateLogic(context.Background(), svcCtx)
	resp, err := l.ChatGroupCreate(&chat.ChatGroupCreateRequest{Name: "测试群组", OperatorUserId: 1})

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

	l := NewChatGroupCreateLogic(context.Background(), svcCtx)
	resp, err := l.ChatGroupCreate(&chat.ChatGroupCreateRequest{Name: "测试群组", OperatorUserId: 1})

	require.Error(t, err)
	assert.Nil(t, resp)
	// 断言：群组那条 INSERT 虽然执行成功过，但整个事务未提交，不会留下没有任何成员的孤儿群组。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
