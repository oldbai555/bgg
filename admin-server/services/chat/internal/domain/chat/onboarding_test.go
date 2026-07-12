package chat

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
	"google.golang.org/grpc"

	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/chat/internal/repository"
)

// fakeIamCallbackClient 是 iamcallbackpb.IamCallbackClient 的最小 fake：内嵌 nil 接口满足
// 全部方法签名，只覆盖测试用到的 FindActiveUserChunk——和
// internal/rpcserver/taskcallback/server_test.go 的 fakeSdkClient 同一个模式（chat-rpc 已经
// 拆成独立进程，不能再用 sqlmock 直接命中 admin_user 表，改成对着这个 RPC 边界打桩）。
type fakeIamCallbackClient struct {
	iamcallbackpb.IamCallbackClient
	findActiveUserChunkFn func(ctx context.Context, in *iamcallbackpb.FindActiveUserChunkRequest) (*iamcallbackpb.FindActiveUserChunkResponse, error)
}

func (f *fakeIamCallbackClient) FindActiveUserChunk(ctx context.Context, in *iamcallbackpb.FindActiveUserChunkRequest, _ ...grpc.CallOption) (*iamcallbackpb.FindActiveUserChunkResponse, error) {
	return f.findActiveUserChunkFn(ctx, in)
}

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

// TestCreatePrivateChat_HappyPath / _Rollback 从 internal/domain/chat/onboarding_test.go
// 原样迁移（chat 域拆分后 store.Transact 换掉了原来的 repo.Transact，方法体本身不变）。
func TestCreatePrivateChat_HappyPath(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
		WillReturnResult(sqlmock.NewResult(9, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(2, 1))
	sqlMock.ExpectCommit()

	s := NewChatOnboardingService(store, nil)
	err := s.createPrivateChat(context.Background(), 1, 2)

	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestCreatePrivateChat_RollbackOnSecondMemberInsertError(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
		WillReturnResult(sqlmock.NewResult(9, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	s := NewChatOnboardingService(store, nil)
	err := s.createPrivateChat(context.Background(), 1, 2)

	require.Error(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

// TestCreatePrivateChatsForExistingUsers_PaginatesViaIamCallback 验证分页边界（page1 恰好
// 100 条触发第二次 FindActiveUserChunk，page2 不足 100 条提前结束）——迁移自原
// internal/domain/chat/onboarding_pagination_test.go，原来打桩的是 chatdomain.UserLister
// mock，现在打桩的是回调单体 IamCallback 的 gRPC 客户端，断言的分页语义完全不变。
func TestCreatePrivateChatsForExistingUsers_PaginatesViaIamCallback(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	// 第一批 100 个用户，最后一个 lastId=100 -> 触发第二次调用
	page1Users := make([]*iamcallbackpb.ActiveUserRef, 100)
	for i := range page1Users {
		page1Users[i] = &iamcallbackpb.ActiveUserRef{Id: uint64(i + 1)}
	}
	// 第二批只有 1 个用户，不足 100 -> 结束分页
	page2Users := []*iamcallbackpb.ActiveUserRef{{Id: 101}}

	callCount := 0
	fakeClient := &fakeIamCallbackClient{
		findActiveUserChunkFn: func(ctx context.Context, in *iamcallbackpb.FindActiveUserChunkRequest) (*iamcallbackpb.FindActiveUserChunkResponse, error) {
			callCount++
			switch callCount {
			case 1:
				assert.Equal(t, uint64(0), in.LastId)
				return &iamcallbackpb.FindActiveUserChunkResponse{Users: page1Users, NextLastId: 100}, nil
			case 2:
				assert.Equal(t, uint64(100), in.LastId)
				return &iamcallbackpb.FindActiveUserChunkResponse{Users: page2Users, NextLastId: 101}, nil
			default:
				t.Fatalf("unexpected call count: %d", callCount)
				return nil, nil
			}
		},
	}

	// newUserID=999 不与任何存量用户重合，101 个存量用户中每一个都要先查是否已有私聊
	// （FindPrivateChatByUserIDs），查不到（sql.ErrNoRows 类空结果集）再建私聊，
	// 每次建私聊是一个独立事务（chat + 2 条 chat_user）。
	sqlMock.MatchExpectationsInOrder(false)
	for i := 0; i < 101; i++ {
		sqlMock.ExpectQuery(regexp.QuoteMeta("FROM chat c")).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		sqlMock.ExpectBegin()
		sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat`")).
			WillReturnResult(sqlmock.NewResult(int64(1000+i), 1))
		sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
			WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectExec(regexp.QuoteMeta("insert into `chat_user`")).
			WillReturnResult(sqlmock.NewResult(1, 1))
		sqlMock.ExpectCommit()
	}

	s := NewChatOnboardingService(store, fakeClient)
	err := s.createPrivateChatsForExistingUsers(context.Background(), 999)

	require.NoError(t, err)
	assert.Equal(t, 2, callCount, "应该恰好调用两次 FindActiveUserChunk（100 条触发第二页，第二页 1 条结束）")
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
