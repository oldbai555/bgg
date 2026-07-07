package blog

import (
	"postapocgame/admin-server/internal/repository"
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	blogmodel "postapocgame/admin-server/internal/model/blog"
)

type BlogArticleAuditRepository interface {
	Create(ctx context.Context, audit *blogmodel.BlogArticleAudit) error
	LatestByArticleID(ctx context.Context, articleID uint64) (*blogmodel.BlogArticleAudit, error)
}

type blogArticleAuditRepository struct {
	model blogmodel.BlogArticleAuditModel
	conn  sqlx.SqlConn
}

func NewBlogArticleAuditRepository(repo *repository.Repository) BlogArticleAuditRepository {
	return &blogArticleAuditRepository{
		model: repo.BlogArticleAuditModel,
		conn:  repo.DB,
	}
}

func (r *blogArticleAuditRepository) Create(ctx context.Context, audit *blogmodel.BlogArticleAudit) error {
	now := time.Now().Unix()
	if audit.CreatedAt == 0 {
		audit.CreatedAt = now
	}
	if audit.UpdatedAt == 0 {
		audit.UpdatedAt = now
	}
	_, err := r.model.Insert(ctx, audit)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "创建审核记录失败", err)
	}
	return nil
}

func (r *blogArticleAuditRepository) LatestByArticleID(ctx context.Context, articleID uint64) (*blogmodel.BlogArticleAudit, error) {
	sql, args, err := sq.Select(
		"id", "article_id", "audit_status", "audit_remark", "auditor_id", "auditor_name", "created_at", "updated_at", "deleted_at",
	).
		From("`blog_article_audit`").
		Where(sq.And{
			sq.Eq{"article_id": articleID},
			sq.Eq{"deleted_at": 0},
		}).
		OrderBy("id DESC").
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "审核记录 SQL 生成失败", err)
	}

	var audit blogmodel.BlogArticleAudit
	if err := r.conn.QueryRowCtx(ctx, &audit, sql, args...); err != nil {
		return nil, err
	}
	return &audit, nil
}
