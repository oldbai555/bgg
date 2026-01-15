package repository

import (
	"context"
	"time"

	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// BlogFriendLinkRepository 友情链接仓储接口
type BlogFriendLinkRepository interface {
	FindPage(ctx context.Context, page, pageSize int64, status int64, keyword string) ([]model.BlogFriendLink, int64, error)
	FindByID(ctx context.Context, id uint64) (*model.BlogFriendLink, error)
	FindEnabledList(ctx context.Context) ([]model.BlogFriendLink, error)
	Create(ctx context.Context, link *model.BlogFriendLink) error
	Update(ctx context.Context, link *model.BlogFriendLink) error
	Delete(ctx context.Context, id uint64) error
}

type blogFriendLinkRepository struct {
	model model.BlogFriendLinkModel
	conn  sqlx.SqlConn
}

func NewBlogFriendLinkRepository(repo *Repository) BlogFriendLinkRepository {
	return &blogFriendLinkRepository{
		model: repo.BlogFriendLinkModel,
		conn:  repo.DB,
	}
}

func (r *blogFriendLinkRepository) FindByID(ctx context.Context, id uint64) (*model.BlogFriendLink, error) {
	return r.model.FindOne(ctx, id)
}

func (r *blogFriendLinkRepository) FindPage(ctx context.Context, page, pageSize int64, status int64, keyword string) ([]model.BlogFriendLink, int64, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
	}

	if status > 0 {
		conditions = append(conditions, sq.Eq{"status": status})
	}

	if keyword != "" {
		k := "%" + keyword + "%"
		conditions = append(conditions, sq.Or{
			sq.Like{"name": k},
			sq.Like{"remark": k},
		})
	}

	// 统计总数
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`blog_friend_link`").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "友情链接列表统计 SQL 生成失败", err)
	}
	var total int64
	if err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "友情链接列表统计查询失败", err)
	}
	if total == 0 {
		return []model.BlogFriendLink{}, 0, nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	listSQL, listArgs, err := sq.Select("*").
		From("`blog_friend_link`").
		Where(conditions).
		OrderBy("order_num ASC, id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "友情链接列表 SQL 生成失败", err)
	}

	var list []model.BlogFriendLink
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "友情链接列表查询失败", err)
	}

	return list, total, nil
}

func (r *blogFriendLinkRepository) FindEnabledList(ctx context.Context) ([]model.BlogFriendLink, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"status": 1}, // 1=启用
	}

	listSQL, listArgs, err := sq.Select("*").
		From("`blog_friend_link`").
		Where(conditions).
		OrderBy("order_num ASC, id DESC").
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "启用友情链接列表 SQL 生成失败", err)
	}

	var list []model.BlogFriendLink
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "启用友情链接列表查询失败", err)
	}

	return list, nil
}

func (r *blogFriendLinkRepository) Create(ctx context.Context, link *model.BlogFriendLink) error {
	now := time.Now().Unix()
	link.CreatedAt = now
	link.UpdatedAt = now
	link.DeletedAt = 0

	_, err := r.model.Insert(ctx, link)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "创建友情链接失败", err)
	}
	return nil
}

func (r *blogFriendLinkRepository) Update(ctx context.Context, link *model.BlogFriendLink) error {
	link.UpdatedAt = time.Now().Unix()
	err := r.model.Update(ctx, link)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新友情链接失败", err)
	}
	return nil
}

func (r *blogFriendLinkRepository) Delete(ctx context.Context, id uint64) error {
	err := r.model.Delete(ctx, id)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "删除友情链接失败", err)
	}
	return nil
}
