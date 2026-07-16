package iam

import (
	"context"
	"postapocgame/admin-server/services/iam/internal/repository"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint64) (*iammodel.AdminUser, error)
	FindByUsername(ctx context.Context, username string) (*iammodel.AdminUser, error)
	FindPage(ctx context.Context, page, pageSize int64, name string) ([]iammodel.AdminUser, int64, error)
	FindChunk(ctx context.Context, limit int64, lastId uint64) ([]iammodel.AdminUser, uint64, error)
	Create(ctx context.Context, user *iammodel.AdminUser) error
	Update(ctx context.Context, user *iammodel.AdminUser) error
	DeleteByID(ctx context.Context, id uint64) error
}

type userRepository struct {
	model iammodel.AdminUserModel
	conn  sqlx.SqlConn
}

func NewUserRepository(repo *repository.Repository) UserRepository {
	return &userRepository{model: repo.AdminUserModel, conn: repo.DB}
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*iammodel.AdminUser, error) {
	return r.model.FindOne(ctx, id)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*iammodel.AdminUser, error) {
	return r.model.FindOneByUsername(ctx, username)
}

// FindPage 支持用户名模糊查询，基于生成的无缓存查询能力。
func (r *userRepository) FindPage(ctx context.Context, page, pageSize int64, name string) ([]iammodel.AdminUser, int64, error) {
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
		list  []iammodel.AdminUser
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

func (r *userRepository) FindChunk(ctx context.Context, limit int64, lastId uint64) ([]iammodel.AdminUser, uint64, error) {
	return r.model.FindChunk(ctx, limit, lastId)
}

func (r *userRepository) Create(ctx context.Context, user *iammodel.AdminUser) error {
	result, err := r.model.Insert(ctx, user)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = uint64(id)
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *iammodel.AdminUser) error {
	return r.model.Update(ctx, user)
}

func (r *userRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
