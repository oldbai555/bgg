package blog

import (
	"postapocgame/admin-server/services/content/internal/repository"
	"context"

	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
)

type BlogArticleTagRepository interface {
	FindTagsByArticleID(ctx context.Context, articleID uint64) ([]blogmodel.BlogTag, error)
	FindTagsByArticleIDs(ctx context.Context, articleIDs []uint64) (map[uint64][]blogmodel.BlogTag, error)
}

type blogArticleTagRepository struct {
	conn sqlx.SqlConn
}

func NewBlogArticleTagRepository(store *repository.Store) BlogArticleTagRepository {
	return &blogArticleTagRepository{conn: store.DB}
}

func (r *blogArticleTagRepository) FindTagsByArticleID(ctx context.Context, articleID uint64) ([]blogmodel.BlogTag, error) {
	m, err := r.FindTagsByArticleIDs(ctx, []uint64{articleID})
	if err != nil {
		return nil, err
	}
	return m[articleID], nil
}

func (r *blogArticleTagRepository) FindTagsByArticleIDs(ctx context.Context, articleIDs []uint64) (map[uint64][]blogmodel.BlogTag, error) {
	result := make(map[uint64][]blogmodel.BlogTag, len(articleIDs))
	if len(articleIDs) == 0 {
		return result, nil
	}

	// 关联表无联合唯一约束，按 deleted_at 过滤
	sql, args, err := sq.Select(
		"bat.article_id",
		"bt.id", "bt.name", "bt.status", "bt.remark", "bt.created_at", "bt.updated_at", "bt.deleted_at",
	).
		From("`blog_article_tag` bat").
		Join("`blog_tag` bt ON bt.id = bat.tag_id").
		Where(sq.And{
			sq.Eq{"bat.deleted_at": 0},
			sq.Eq{"bt.deleted_at": 0},
			sq.Eq{"bat.article_id": articleIDs},
		}).
		OrderBy("bat.article_id ASC", "bt.id ASC").
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "文章标签查询 SQL 生成失败", err)
	}

	type row struct {
		ArticleId uint64 `db:"article_id"`
		blogmodel.BlogTag
	}
	var rows []row
	if err = r.conn.QueryRowsCtx(ctx, &rows, sql, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "文章标签查询失败", err)
	}

	for _, it := range rows {
		result[it.ArticleId] = append(result[it.ArticleId], it.BlogTag)
	}
	return result, nil
}
