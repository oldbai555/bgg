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

// BlogFriendLinkRepository 友情链接仓储接口
type BlogFriendLinkRepository interface {
	FindPage(ctx context.Context, page, pageSize int64, status int64, keyword string) ([]blogmodel.BlogFriendLink, int64, error)
	FindByID(ctx context.Context, id uint64) (*blogmodel.BlogFriendLink, error)
	FindEnabledList(ctx context.Context) ([]blogmodel.BlogFriendLink, error)
	Create(ctx context.Context, link *blogmodel.BlogFriendLink) error
	Update(ctx context.Context, link *blogmodel.BlogFriendLink) error
	Delete(ctx context.Context, id uint64) error
}

type blogFriendLinkRepository struct {
	model blogmodel.BlogFriendLinkModel
	conn  sqlx.SqlConn
}

func NewBlogFriendLinkRepository(repo *repository.Repository) BlogFriendLinkRepository {
	return &blogFriendLinkRepository{
		model: repo.BlogFriendLinkModel,
		conn:  repo.DB,
	}
}

func (r *blogFriendLinkRepository) FindByID(ctx context.Context, id uint64) (*blogmodel.BlogFriendLink, error) {
	return r.model.FindOne(ctx, id)
}

func (r *blogFriendLinkRepository) FindPage(ctx context.Context, page, pageSize int64, status int64, keyword string) ([]blogmodel.BlogFriendLink, int64, error) {
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
		return []blogmodel.BlogFriendLink{}, 0, nil
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

	var list []blogmodel.BlogFriendLink
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "友情链接列表查询失败", err)
	}

	return list, total, nil
}

func (r *blogFriendLinkRepository) FindEnabledList(ctx context.Context) ([]blogmodel.BlogFriendLink, error) {
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

	var list []blogmodel.BlogFriendLink
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "启用友情链接列表查询失败", err)
	}

	return list, nil
}

func (r *blogFriendLinkRepository) Create(ctx context.Context, link *blogmodel.BlogFriendLink) error {
	now := time.Now().Unix()
	link.CreatedAt = now
	link.UpdatedAt = now
	link.DeletedAt = 0

	result, err := r.model.Insert(ctx, link)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "创建友情链接失败", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "获取友情链接自增 ID 失败", err)
	}
	link.Id = uint64(id)
	return nil
}

func (r *blogFriendLinkRepository) Update(ctx context.Context, link *blogmodel.BlogFriendLink) error {
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
