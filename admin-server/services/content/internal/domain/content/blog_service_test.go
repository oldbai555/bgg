package content_test

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

	"postapocgame/admin-server/pkg/errs"
	contentdomain "postapocgame/admin-server/services/content/internal/domain/content"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
	"postapocgame/admin-server/services/content/internal/repository"
	"postapocgame/admin-server/services/content/internal/consts"
)

var blogArticleColumns = []string{
	"id", "title", "content", "status", "audit_status", "cover", "author_id", "author_name",
	"publish_time", "summary", "is_top", "created_at", "updated_at", "deleted_at",
}

// newTestStore 和 services/sdk/internal/domain/sdk/sdk_service_test.go 用的是同一套组合：
// miniredis + 单节点 CacheConf——goctl 生成的 Model 内部走 sqlc.CachedConn，cache.New 对空
// CacheConf 会 log.Fatal，必须喂一个真实（哪怕是内存模拟的）Redis 节点。
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

func TestBlogArticleService_CreateArticle_HappyPath(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article_tag`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article_tag`")).
		WillReturnResult(sqlmock.NewResult(2, 1))
	sqlMock.ExpectCommit()

	svc := contentdomain.NewBlogArticleService(store)
	article := &blogmodel.BlogArticle{Title: "标题", Content: "内容"}
	err := svc.CreateArticle(context.Background(), article, []uint64{10, 11})

	require.NoError(t, err)
	assert.Equal(t, uint64(1), article.Id)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_CreateArticle_RollbackOnTagInsertError(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article_tag`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `blog_article_tag`")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	svc := contentdomain.NewBlogArticleService(store)
	article := &blogmodel.BlogArticle{Title: "标题", Content: "内容"}
	err := svc.CreateArticle(context.Background(), article, []uint64{10, 11})

	require.Error(t, err)
	// 事务整体回滚，article 那条 INSERT 不会被提交，也不需要任何手动补偿删除逻辑。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_AuditArticle_HappyPath(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(1, "标题", "内容", consts.BlogArticleStatusPendingAudit, consts.BlogArticleAuditStatusPending, "", 1, "author", 0, "", 0, 0, 0, 0))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `blog_article_audit`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("update `blog_article`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectCommit()

	svc := contentdomain.NewBlogArticleService(store)
	article, err := svc.AuditArticle(context.Background(), 1, consts.BlogArticleAuditStatusPassed, "通过", 2, "auditor")

	require.NoError(t, err)
	require.NotNil(t, article)
	assert.Equal(t, consts.BlogArticleAuditStatusPassed, article.AuditStatus)
	assert.Equal(t, consts.BlogArticleStatusAuditPassed, article.Status)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_AuditArticle_RollbackOnNotPending(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(1, "标题", "内容", consts.BlogArticleStatusAuditPassed, consts.BlogArticleAuditStatusPassed, "", 1, "author", 0, "", 0, 0, 0, 0))
	sqlMock.ExpectRollback()

	svc := contentdomain.NewBlogArticleService(store)
	article, err := svc.AuditArticle(context.Background(), 1, consts.BlogArticleAuditStatusPassed, "通过", 2, "auditor")

	require.Error(t, err)
	assert.Nil(t, article)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeForbidden, bizErr.Code)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_AuditArticle_RollbackOnArticleUpdateError(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(1, "标题", "内容", consts.BlogArticleStatusPendingAudit, consts.BlogArticleAuditStatusPending, "", 1, "author", 0, "", 0, 0, 0, 0))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `blog_article_audit`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("update `blog_article`")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	svc := contentdomain.NewBlogArticleService(store)
	article, err := svc.AuditArticle(context.Background(), 1, consts.BlogArticleAuditStatusPassed, "通过", 2, "auditor")

	require.Error(t, err)
	assert.Nil(t, article)
	// 审核记录那条 INSERT 虽然执行成功过，但整个事务未提交，最终会随 Rollback 一起撤销。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_UnpublishArticle_HappyPath(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(1, "标题", "内容", consts.BlogArticleStatusPublished, consts.BlogArticleAuditStatusPassed, "", 1, "author", 0, "", 0, 0, 0, 0))
	sqlMock.ExpectExec(regexp.QuoteMeta("update `blog_article`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(regexp.QuoteMeta("insert into `blog_article_audit`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	sqlMock.ExpectCommit()

	svc := contentdomain.NewBlogArticleService(store)
	article, err := svc.UnpublishArticle(context.Background(), 1, "下架原因", 2, "operator")

	require.NoError(t, err)
	require.NotNil(t, article)
	assert.Equal(t, consts.BlogArticleStatusUnpublished, article.Status)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_SetArticleTop_HappyPath_CancelsOldest(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	// FindByID(2)：目标文章未置顶
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(2, "标题2", "内容", consts.BlogArticleStatusPublished, consts.BlogArticleAuditStatusPassed, "", 1, "author", 0, "", 0, 0, 0, 0))
	// FindTopCount：已达上限（squirrel 生成大写 SELECT，和 goctl Model 的小写 select 要分开匹配）
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `blog_article`")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	// FindOldestTopArticle：id=1
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(1, "标题1", "内容", consts.BlogArticleStatusPublished, consts.BlogArticleAuditStatusPassed, "", 1, "author", 0, "", 0, 1, 0, 0))
	// UpdateTopStatus(1, 0)：取消最早置顶 + 内部刷新缓存的 FindOne
	sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `blog_article`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(1, "标题1", "内容", consts.BlogArticleStatusPublished, consts.BlogArticleAuditStatusPassed, "", 1, "author", 0, "", 0, 0, 0, 0))
	// UpdateTopStatus(2, 1)：置顶目标文章。内部刷新缓存的 FindOne(2) 命中的是 FindByID(2)
	// 在事务开头写入的缓存（同一个 miniredis 实例，未做事务级隔离），不会再打到 DB。
	sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `blog_article`")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectCommit()

	svc := contentdomain.NewBlogArticleService(store)
	err := svc.SetArticleTop(context.Background(), 2, 1)

	require.NoError(t, err)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_SetArticleTop_RollbackOnFinalUpdateError(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	// FindByID(2)：目标文章未置顶
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(2, "标题2", "内容", consts.BlogArticleStatusPublished, consts.BlogArticleAuditStatusPassed, "", 1, "author", 0, "", 0, 0, 0, 0))
	// FindTopCount：未达上限，跳过取消最早置顶分支
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM `blog_article`")).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// UpdateTopStatus(2, 1) 失败
	sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `blog_article`")).
		WillReturnError(assert.AnError)
	sqlMock.ExpectRollback()

	svc := contentdomain.NewBlogArticleService(store)
	err := svc.SetArticleTop(context.Background(), 2, 1)

	require.Error(t, err)
	// 断言：置顶数量查询等只读操作不需要任何补偿，整个事务随 Rollback 一起撤销。
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestBlogArticleService_UnpublishArticle_RollbackOnNotPublished(t *testing.T) {
	store, sqlMock, cleanup := newTestStore(t)
	defer cleanup()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("from `blog_article`")).
		WillReturnRows(sqlmock.NewRows(blogArticleColumns).
			AddRow(1, "标题", "内容", consts.BlogArticleStatusDraft, consts.BlogArticleAuditStatusNotSubmitted, "", 1, "author", 0, "", 0, 0, 0, 0))
	sqlMock.ExpectRollback()

	svc := contentdomain.NewBlogArticleService(store)
	article, err := svc.UnpublishArticle(context.Background(), 1, "下架原因", 2, "operator")

	require.Error(t, err)
	assert.Nil(t, article)
	bizErr, ok := errs.FromError(err)
	require.True(t, ok)
	assert.Equal(t, errs.CodeForbidden, bizErr.Code)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}
