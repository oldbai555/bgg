package blog

import (
	"postapocgame/admin-server/services/content/internal/repository"
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
)

// BlogSocialInfoRepository 社交信息仓储接口
type BlogSocialInfoRepository interface {
	FindPage(ctx context.Context, page, pageSize int64, status int64, keyword string) ([]blogmodel.BlogSocialInfo, int64, error)
	FindByID(ctx context.Context, id uint64) (*blogmodel.BlogSocialInfo, error)
	FindEnabledList(ctx context.Context) ([]blogmodel.BlogSocialInfo, error)
	Create(ctx context.Context, info *blogmodel.BlogSocialInfo) error
	Update(ctx context.Context, info *blogmodel.BlogSocialInfo) error
	Delete(ctx context.Context, id uint64) error
}

type blogSocialInfoRepository struct {
	model blogmodel.BlogSocialInfoModel
	conn  sqlx.SqlConn
}

func NewBlogSocialInfoRepository(store *repository.Store) BlogSocialInfoRepository {
	return &blogSocialInfoRepository{
		model: store.BlogSocialInfoModel,
		conn:  store.DB,
	}
}

func (r *blogSocialInfoRepository) FindByID(ctx context.Context, id uint64) (*blogmodel.BlogSocialInfo, error) {
	return r.model.FindOne(ctx, id)
}

func (r *blogSocialInfoRepository) FindPage(ctx context.Context, page, pageSize int64, status int64, keyword string) ([]blogmodel.BlogSocialInfo, int64, error) {
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
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`blog_social_info`").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "社交信息列表统计 SQL 生成失败", err)
	}
	var total int64
	if err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "社交信息列表统计查询失败", err)
	}
	if total == 0 {
		return []blogmodel.BlogSocialInfo{}, 0, nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	listSQL, listArgs, err := sq.Select("*").
		From("`blog_social_info`").
		Where(conditions).
		OrderBy("order_num ASC, id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "社交信息列表 SQL 生成失败", err)
	}

	var list []blogmodel.BlogSocialInfo
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "社交信息列表查询失败", err)
	}

	return list, total, nil
}

func (r *blogSocialInfoRepository) FindEnabledList(ctx context.Context) ([]blogmodel.BlogSocialInfo, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"status": 1}, // 1=启用
	}

	listSQL, listArgs, err := sq.Select("*").
		From("`blog_social_info`").
		Where(conditions).
		OrderBy("order_num ASC, id DESC").
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "启用社交信息列表 SQL 生成失败", err)
	}

	var list []blogmodel.BlogSocialInfo
	if err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "启用社交信息列表查询失败", err)
	}

	return list, nil
}

func (r *blogSocialInfoRepository) Create(ctx context.Context, info *blogmodel.BlogSocialInfo) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	info.DeletedAt = 0

	result, err := r.model.Insert(ctx, info)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "创建社交信息失败", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "获取社交信息自增 ID 失败", err)
	}
	info.Id = uint64(id)
	return nil
}

func (r *blogSocialInfoRepository) Update(ctx context.Context, info *blogmodel.BlogSocialInfo) error {
	info.UpdatedAt = time.Now().Unix()
	err := r.model.Update(ctx, info)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "更新社交信息失败", err)
	}
	return nil
}

func (r *blogSocialInfoRepository) Delete(ctx context.Context, id uint64) error {
	err := r.model.Delete(ctx, id)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "删除社交信息失败", err)
	}
	return nil
}
