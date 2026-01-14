package repository

import (
	"context"

	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// BlogTagRepository 博客标签仓储接口
type BlogTagRepository interface {
	FindPage(ctx context.Context, page, pageSize int64, name string, status int64) ([]model.BlogTag, int64, error)
	FindByID(ctx context.Context, id uint64) (*model.BlogTag, error)
	// FindEnabledList 查询启用标签列表（用于下拉选项）
	FindEnabledList(ctx context.Context, limit int64) ([]model.BlogTag, error)
	Create(ctx context.Context, tag *model.BlogTag) error
	Update(ctx context.Context, tag *model.BlogTag) error
	Delete(ctx context.Context, id uint64) error
}

type blogTagRepository struct {
	model model.BlogTagModel
	conn  sqlx.SqlConn
}

func NewBlogTagRepository(repo *Repository) BlogTagRepository {
	return &blogTagRepository{
		model: repo.BlogTagModel,
		conn:  repo.DB,
	}
}

func (r *blogTagRepository) FindByID(ctx context.Context, id uint64) (*model.BlogTag, error) {
	return r.model.FindOne(ctx, id)
}

func (r *blogTagRepository) Create(ctx context.Context, tag *model.BlogTag) error {
	_, err := r.model.Insert(ctx, tag)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "创建标签失败", err)
	}
	return nil
}

func (r *blogTagRepository) Update(ctx context.Context, tag *model.BlogTag) error {
	err := r.model.Update(ctx, tag)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新标签失败", err)
	}
	return nil
}

func (r *blogTagRepository) Delete(ctx context.Context, id uint64) error {
	err := r.model.Delete(ctx, id)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "删除标签失败", err)
	}
	return nil
}

func (r *blogTagRepository) FindPage(ctx context.Context, page, pageSize int64, name string, status int64) ([]model.BlogTag, int64, error) {
	// 使用 squirrel 构建动态查询
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
	}

	if name != "" {
		conditions = append(conditions, sq.Like{"name": "%" + name + "%"})
	}

	// status 字段使用字典值，0 表示不筛选
	if status > 0 {
		conditions = append(conditions, sq.Eq{"status": status})
	}

	// 统计总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").
		From("`blog_tag`").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "标签列表统计 SQL 生成失败", err)
	}
	if err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "标签列表统计查询失败", err)
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// 查询列表
	listSQL, listArgs, err := sq.Select("id", "name", "status", "remark", "created_at", "updated_at", "deleted_at").
		From("`blog_tag`").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "标签列表 SQL 生成失败", err)
	}

	var list []model.BlogTag
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "标签列表查询失败", err)
	}

	return list, total, nil
}

func (r *blogTagRepository) FindEnabledList(ctx context.Context, limit int64) ([]model.BlogTag, error) {
	if limit <= 0 {
		limit = 1000
	}
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"status": 1},
	}
	// 查询所有字段，以匹配 model.BlogTag 结构体
	sql, args, err := sq.Select("id", "name", "status", "remark", "created_at", "updated_at", "deleted_at").
		From("`blog_tag`").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "标签选项 SQL 生成失败", err)
	}
	var list []model.BlogTag
	if err = r.conn.QueryRowsCtx(ctx, &list, sql, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "标签选项查询失败", err)
	}
	return list, nil
}
