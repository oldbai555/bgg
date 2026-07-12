package logic_test

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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/config"
	contentdomain "postapocgame/admin-server/services/content/internal/domain/content"
	"postapocgame/admin-server/services/content/internal/logic"
	"postapocgame/admin-server/services/content/internal/repository"
	blogrepo "postapocgame/admin-server/services/content/internal/repository/blog"
	videorepo "postapocgame/admin-server/services/content/internal/repository/video"
	"postapocgame/admin-server/services/content/internal/svc"
)

// fakeIamCallbackClient 是 iamcallbackpb.IamCallbackClient 的最小 fake：内嵌 nil 接口满足
// 全部方法签名，只覆盖测试用到的 GetUserProfile/RecordAuditLog——和
// services/chat/internal/domain/chat/onboarding_test.go 的 fakeIamCallbackClient 同一个模式。
type fakeIamCallbackClient struct {
	iamcallbackpb.IamCallbackClient
	getUserProfileFn  func(ctx context.Context, in *iamcallbackpb.GetUserProfileRequest) (*iamcallbackpb.GetUserProfileResponse, error)
	recordAuditLogFn  func(ctx context.Context, in *iamcallbackpb.RecordAuditLogRequest) (*iamcallbackpb.RecordAuditLogResponse, error)
}

func (f *fakeIamCallbackClient) GetUserProfile(ctx context.Context, in *iamcallbackpb.GetUserProfileRequest, _ ...grpc.CallOption) (*iamcallbackpb.GetUserProfileResponse, error) {
	return f.getUserProfileFn(ctx, in)
}

func (f *fakeIamCallbackClient) RecordAuditLog(ctx context.Context, in *iamcallbackpb.RecordAuditLogRequest, _ ...grpc.CallOption) (*iamcallbackpb.RecordAuditLogResponse, error) {
	if f.recordAuditLogFn == nil {
		return &iamcallbackpb.RecordAuditLogResponse{}, nil
	}
	return f.recordAuditLogFn(ctx, in)
}

// newTestSvcCtx 和 domain 包的 newTestStore 用的是同一套组合（miniredis + 单节点 CacheConf，
// cache.New 对空 CacheConf 会 log.Fatal），直接拼出 logic 需要的完整 *svc.ServiceContext。
func newTestSvcCtx(t *testing.T, iamCallback iamcallbackpb.IamCallbackClient) (*svc.ServiceContext, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	conn := sqlx.NewSqlConnFromDB(db)

	mr, err := miniredis.Run()
	require.NoError(t, err)

	redisConf := redis.RedisConf{Host: mr.Addr(), Type: "node"}
	cacheConf := cache.CacheConf{{RedisConf: redisConf, Weight: 100}}

	store := repository.NewStore(conn, cacheConf)

	cfg := config.Config{}
	cfg.Limits.BlogTagNameMaxLength = 10
	cfg.Limits.BlogArticleTitleMaxLength = 100
	cfg.Limits.BlogArticleSummaryLength = 120
	cfg.Limits.BlogArticleTopMaxCount = 1
	cfg.Limits.BlogFriendLinkNameMaxLength = 15
	cfg.Limits.BlogFriendLinkUrlMaxLength = 255
	cfg.Limits.BlogFriendLinkRemarkMaxLength = 127
	cfg.Limits.BlogSocialInfoNameMaxLength = 15
	cfg.Limits.BlogSocialInfoUrlMaxLength = 255
	cfg.Limits.BlogSocialInfoRemarkMaxLength = 127

	svcCtx := &svc.ServiceContext{
		Config:         cfg,
		Store:          store,
		BlogArticle:    blogrepo.NewBlogArticleRepository(store),
		BlogArticleTag: blogrepo.NewBlogArticleTagRepository(store),
		ArticleAudit:   blogrepo.NewBlogArticleAuditRepository(store),
		FriendLink:     blogrepo.NewBlogFriendLinkRepository(store),
		SocialInfo:     blogrepo.NewBlogSocialInfoRepository(store),
		Tag:            blogrepo.NewBlogTagRepository(store),
		Video:          videorepo.NewVideoRepository(store),
		ArticleService: contentdomain.NewBlogArticleService(store),
		IamCallback:    iamCallback,
	}

	return svcCtx, sqlMock, func() {
		_ = db.Close()
		mr.Close()
	}
}

func grpcCode(t *testing.T, err error) codes.Code {
	t.Helper()
	st, ok := status.FromError(err)
	require.True(t, ok, "expected a gRPC status error, got %v", err)
	return st.Code()
}

func TestBlogTagCreate_NameTooLong(t *testing.T) {
	svcCtx, _, cleanup := newTestSvcCtx(t, nil)
	defer cleanup()

	_, err := logic.NewBlogTagCreateLogic(context.Background(), svcCtx).BlogTagCreate(&content.BlogTagCreateRequest{
		Name: "这是一个超过十个字符的标签名称",
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, grpcCode(t, err))
}

func TestBlogTagCreate_NameEmpty(t *testing.T) {
	svcCtx, _, cleanup := newTestSvcCtx(t, nil)
	defer cleanup()

	_, err := logic.NewBlogTagCreateLogic(context.Background(), svcCtx).BlogTagCreate(&content.BlogTagCreateRequest{Name: "  "})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, grpcCode(t, err))
}

func TestBlogArticleCreate_TitleTooLong(t *testing.T) {
	svcCtx, _, cleanup := newTestSvcCtx(t, nil)
	defer cleanup()

	longTitle := ""
	for i := 0; i < 101; i++ {
		longTitle += "字"
	}

	_, err := logic.NewBlogArticleCreateLogic(context.Background(), svcCtx).BlogArticleCreate(&content.BlogArticleCreateRequest{
		Title:  longTitle,
		TagIds: []uint64{1},
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, grpcCode(t, err))
}

func TestBlogArticleCreate_NoTags(t *testing.T) {
	svcCtx, _, cleanup := newTestSvcCtx(t, nil)
	defer cleanup()

	_, err := logic.NewBlogArticleCreateLogic(context.Background(), svcCtx).BlogArticleCreate(&content.BlogArticleCreateRequest{
		Title: "标题",
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, grpcCode(t, err))
}

func TestBlogArticleCreate_HappyPath(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, nil)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article_tag`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()

	_, err := logic.NewBlogArticleCreateLogic(context.Background(), svcCtx).BlogArticleCreate(&content.BlogArticleCreateRequest{
		Title:            "标题",
		Content:          "内容",
		TagIds:           []uint64{10},
		OperatorUserId:   2,
		OperatorUsername: "author",
	})
	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleAudit_InvalidResult(t *testing.T) {
	svcCtx, _, cleanup := newTestSvcCtx(t, nil)
	defer cleanup()

	_, err := logic.NewBlogArticleAuditLogic(context.Background(), svcCtx).BlogArticleAudit(&content.BlogArticleAuditRequest{
		Id:     1,
		Result: 99,
	})
	require.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, grpcCode(t, err))
}

func TestPublicBlogAuthorInfo_FallbackWhenNotExists(t *testing.T) {
	fakeClient := &fakeIamCallbackClient{
		getUserProfileFn: func(ctx context.Context, in *iamcallbackpb.GetUserProfileRequest) (*iamcallbackpb.GetUserProfileResponse, error) {
			return &iamcallbackpb.GetUserProfileResponse{Exists: false}, nil
		},
	}
	svcCtx, _, cleanup := newTestSvcCtx(t, fakeClient)
	defer cleanup()

	resp, err := logic.NewPublicBlogAuthorInfoLogic(context.Background(), svcCtx).PublicBlogAuthorInfo(&content.PublicBlogGlobalRequest{})
	require.NoError(t, err)
	assert.Equal(t, uint64(1), resp.Id)
	assert.Equal(t, "管理员", resp.Nickname)
}

func TestPublicBlogAuthorInfo_Success(t *testing.T) {
	fakeClient := &fakeIamCallbackClient{
		getUserProfileFn: func(ctx context.Context, in *iamcallbackpb.GetUserProfileRequest) (*iamcallbackpb.GetUserProfileResponse, error) {
			assert.Equal(t, uint64(1), in.UserId)
			return &iamcallbackpb.GetUserProfileResponse{
				Exists: true, Nickname: "站长", Avatar: "https://a.png", Signature: "hello",
			}, nil
		},
	}
	svcCtx, _, cleanup := newTestSvcCtx(t, fakeClient)
	defer cleanup()

	resp, err := logic.NewPublicBlogAuthorInfoLogic(context.Background(), svcCtx).PublicBlogAuthorInfo(&content.PublicBlogGlobalRequest{})
	require.NoError(t, err)
	assert.Equal(t, "站长", resp.Nickname)
	assert.Equal(t, "https://a.png", resp.Avatar)
	assert.Equal(t, "hello", resp.Signature)
}

func TestVideoCollect_DuplicateUuid(t *testing.T) {
	svcCtx, sqlMock, cleanup := newTestSvcCtx(t, nil)
	defer cleanup()

	videoColumns := []string{"id", "uuid", "name", "cover", "god_num", "duration", "play_url", "xlzz_urls", "description", "type", "created_at", "updated_at", "deleted_at"}
	sqlMock.ExpectQuery(regexp.QuoteMeta("FROM `video`")).
		WillReturnRows(sqlmock.NewRows(videoColumns).
			AddRow(1, "dup-uuid", "existing", "", "", 0, "http://x", "", "", 2, 0, 0, 0))

	_, err := logic.NewVideoCollectLogic(context.Background(), svcCtx).VideoCollect(&content.VideoCollectRequest{
		Uuid:      "dup-uuid",
		PlayerUrl: "http://x",
		Name:      "test",
	})
	require.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, grpcCode(t, err))
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
