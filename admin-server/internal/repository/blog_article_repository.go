package repository

import (
	"context"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// BlogArticleRepository 博客文章仓储接口
type BlogArticleRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.BlogArticle, error)
	FindPage(ctx context.Context, page, pageSize int64, title string, status, auditStatus int64, tagId uint64, startTime, endTime int64) ([]model.BlogArticle, int64, error)
	FindPublicPage(ctx context.Context, page, pageSize int64, keyword string, tagId uint64) ([]model.BlogArticle, int64, error)
	CreateWithTags(ctx context.Context, article *model.BlogArticle, tagIDs []uint64) error
	UpdateWithTags(ctx context.Context, article *model.BlogArticle, tagIDs []uint64) error
	Delete(ctx context.Context, id uint64) error
}

type blogArticleRepository struct {
	articleModel model.BlogArticleModel
	articleTag   model.BlogArticleTagModel
	conn         sqlx.SqlConn
}

func NewBlogArticleRepository(repo *Repository) BlogArticleRepository {
	return &blogArticleRepository{
		articleModel: repo.BlogArticleModel,
		articleTag:   repo.BlogArticleTagModel,
		conn:         repo.DB,
	}
}

func (r *blogArticleRepository) FindByID(ctx context.Context, id uint64) (*model.BlogArticle, error) {
	return r.articleModel.FindOne(ctx, id)
}

func (r *blogArticleRepository) FindPage(ctx context.Context, page, pageSize int64, title string, status, auditStatus int64, tagId uint64, startTime, endTime int64) ([]model.BlogArticle, int64, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
	}

	if title != "" {
		conditions = append(conditions, sq.Like{"title": "%" + title + "%"})
	}
	if status > 0 {
		conditions = append(conditions, sq.Eq{"status": status})
	}
	if auditStatus > 0 {
		conditions = append(conditions, sq.Eq{"audit_status": auditStatus})
	}
	if startTime > 0 {
		conditions = append(conditions, sq.GtOrEq{"created_at": startTime})
	}
	if endTime > 0 {
		conditions = append(conditions, sq.LtOrEq{"created_at": endTime})
	}

	// 标签筛选：使用 EXISTS 子查询（避免 join 导致 count/分页重复）
	if tagId > 0 {
		conditions = append(conditions, sq.Expr(
			"EXISTS (SELECT 1 FROM `blog_article_tag` bat WHERE bat.article_id = `blog_article`.id AND bat.tag_id = ? AND bat.deleted_at = 0)",
			tagId,
		))
	}

	// 统计总数
	countSQL, countArgs, err := sq.Select("COUNT(*)").
		From("`blog_article`").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "文章列表统计 SQL 生成失败", err)
	}

	var total int64
	if err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "文章列表统计查询失败", err)
	}

	if total == 0 {
		return []model.BlogArticle{}, 0, nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 查询列表
	listSQL, listArgs, err := sq.Select("id", "title", "content", "status", "audit_status", "cover", "author_id", "author_name", "publish_time", "summary", "created_at", "updated_at", "deleted_at").
		From("`blog_article`").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "文章列表 SQL 生成失败", err)
	}

	var list []model.BlogArticle
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "文章列表查询失败", err)
	}

	return list, total, nil
}

func (r *blogArticleRepository) FindPublicPage(ctx context.Context, page, pageSize int64, keyword string, tagId uint64) ([]model.BlogArticle, int64, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		// 仅展示已审核通过 + 上架
		sq.Eq{"audit_status": consts.BlogArticleAuditStatusPassed},
		sq.Eq{"status": consts.BlogArticleStatusPublished},
	}

	if keyword != "" {
		k := "%" + keyword + "%"
		conditions = append(conditions, sq.Or{
			sq.Like{"title": k},
			sq.Like{"summary": k},
		})
	}

	if tagId > 0 {
		conditions = append(conditions, sq.Expr(
			"EXISTS (SELECT 1 FROM `blog_article_tag` bat WHERE bat.article_id = `blog_article`.id AND bat.tag_id = ? AND bat.deleted_at = 0)",
			tagId,
		))
	}

	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`blog_article`").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "公共文章列表统计 SQL 生成失败", err)
	}
	var total int64
	if err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "公共文章列表统计查询失败", err)
	}
	if total == 0 {
		return []model.BlogArticle{}, 0, nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	listSQL, listArgs, err := sq.Select("id", "title", "content", "status", "audit_status", "cover", "author_id", "author_name", "publish_time", "summary", "created_at", "updated_at", "deleted_at").
		From("`blog_article`").
		Where(conditions).
		OrderBy("publish_time DESC", "id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "公共文章列表 SQL 生成失败", err)
	}

	var list []model.BlogArticle
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "公共文章列表查询失败", err)
	}

	return list, total, nil
}

func (r *blogArticleRepository) CreateWithTags(ctx context.Context, article *model.BlogArticle, tagIDs []uint64) error {
	now := time.Now().Unix()
	if article.CreatedAt == 0 {
		article.CreatedAt = now
	}
	if article.UpdatedAt == 0 {
		article.UpdatedAt = now
	}

	// 使用 squirrel 手动插入，避免依赖事务 session API
	insertSQL, insertArgs, err := sq.Insert("`blog_article`").
		Columns("`title`", "`content`", "`status`", "`audit_status`", "`cover`", "`author_id`", "`author_name`", "`publish_time`", "`summary`", "`created_at`", "`updated_at`", "`deleted_at`").
		Values(article.Title, article.Content, article.Status, article.AuditStatus, article.Cover, article.AuthorId, article.AuthorName, article.PublishTime, article.Summary, article.CreatedAt, article.UpdatedAt, article.DeletedAt).
		ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "创建文章 SQL 生成失败", err)
	}
	res, err := r.conn.ExecCtx(ctx, insertSQL, insertArgs...)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "创建文章失败", err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "获取文章ID失败", err)
	}
	article.Id = uint64(lastID)

	// 插入标签关联
	for _, tagID := range tagIDs {
		relSQL, relArgs, err := sq.Insert("`blog_article_tag`").
			Columns("`article_id`", "`tag_id`", "`created_at`", "`updated_at`", "`deleted_at`").
			Values(article.Id, tagID, now, now, 0).
			ToSql()
		if err != nil {
			_ = r.articleModel.Delete(ctx, article.Id)
			return errs.Wrap(errs.CodeBadDB, "创建文章标签关联 SQL 生成失败", err)
		}
		if _, err = r.conn.ExecCtx(ctx, relSQL, relArgs...); err != nil {
			_ = r.articleModel.Delete(ctx, article.Id)
			return errs.Wrap(errs.CodeBadDB, "创建文章标签关联失败", err)
		}
	}

	return nil
}

func (r *blogArticleRepository) UpdateWithTags(ctx context.Context, article *model.BlogArticle, tagIDs []uint64) error {
	article.UpdatedAt = time.Now().Unix()
	if err := r.articleModel.Update(ctx, article); err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新文章失败", err)
	}

	// 先删除旧的标签关联（软删除）
	deleteSQL, deleteArgs, err := sq.Update("`blog_article_tag`").
		Set("deleted_at", time.Now().Unix()).
		Where(sq.Eq{"article_id": article.Id, "deleted_at": 0}).
		ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "文章标签删除 SQL 生成失败", err)
	}
	if _, err = r.conn.ExecCtx(ctx, deleteSQL, deleteArgs...); err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新文章标签关联失败", err)
	}

	now := time.Now().Unix()
	for _, tagID := range tagIDs {
		relSQL, relArgs, err := sq.Insert("`blog_article_tag`").
			Columns("`article_id`", "`tag_id`", "`created_at`", "`updated_at`", "`deleted_at`").
			Values(article.Id, tagID, now, now, 0).
			ToSql()
		if err != nil {
			return errs.Wrap(errs.CodeBadDB, "创建文章标签关联 SQL 生成失败", err)
		}
		if _, err = r.conn.ExecCtx(ctx, relSQL, relArgs...); err != nil {
			return errs.Wrap(errs.CodeBadDB, "创建文章标签关联失败", err)
		}
	}

	return nil
}

func (r *blogArticleRepository) Delete(ctx context.Context, id uint64) error {
	if err := r.articleModel.Delete(ctx, id); err != nil {
		return errs.Wrap(errs.CodeBadDB, "删除文章失败", err)
	}
	// 同步软删标签关联
	sql, args, err := sq.Update("`blog_article_tag`").
		Set("deleted_at", time.Now().Unix()).
		Where(sq.Eq{"article_id": id, "deleted_at": 0}).
		ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "文章标签删除 SQL 生成失败", err)
	}
	if _, err = r.conn.ExecCtx(ctx, sql, args...); err != nil {
		return errs.Wrap(errs.CodeBadDB, "删除文章标签关联失败", err)
	}
	return nil
}
