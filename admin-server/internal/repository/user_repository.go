package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.AdminUser, error)
	FindByUsername(ctx context.Context, username string) (*model.AdminUser, error)
	FindPage(ctx context.Context, page, pageSize int64, name string) ([]model.AdminUser, int64, error)
	FindChunk(ctx context.Context, limit int64, lastId uint64) ([]model.AdminUser, uint64, error)
	Create(ctx context.Context, user *model.AdminUser) error
	Update(ctx context.Context, user *model.AdminUser) error
	DeleteByID(ctx context.Context, id uint64) error
}

type userRepository struct {
	model model.AdminUserModel
	conn  sqlx.SqlConn
}

func NewUserRepository(repo *Repository) UserRepository {
	return &userRepository{model: repo.AdminUserModel, conn: repo.DB}
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*model.AdminUser, error) {
	return r.model.FindOne(ctx, id)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.AdminUser, error) {
	return r.model.FindOneByUsername(ctx, username)
}

// FindPage 支持用户名模糊查询，基于生成的无缓存查询能力。
func (r *userRepository) FindPage(ctx context.Context, page, pageSize int64, name string) ([]model.AdminUser, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var (
		list  []model.AdminUser
		total int64
	)

	if name == "" {
		return r.model.FindPage(ctx, page, pageSize)
	}

	// 带用户名模糊筛选的自定义查询
	countSQL, countArgs, err := sq.Select("count(*)").
		From("admin_user").
		Where(sq.And{
			sq.Eq{"deleted_at": 0},
			sq.Like{"username": "%" + name + "%"},
		}).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}
	listSQL, listArgs, err := sq.Select("*").
		From("admin_user").
		Where(sq.And{
			sq.Eq{"deleted_at": 0},
			sq.Like{"username": "%" + name + "%"},
		}).
		OrderBy("id desc").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *userRepository) FindChunk(ctx context.Context, limit int64, lastId uint64) ([]model.AdminUser, uint64, error) {
	return r.model.FindChunk(ctx, limit, lastId)
}

func (r *userRepository) Create(ctx context.Context, user *model.AdminUser) error {
	_, err := r.model.Insert(ctx, user)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *model.AdminUser) error {
	return r.model.Update(ctx, user)
}

func (r *userRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
