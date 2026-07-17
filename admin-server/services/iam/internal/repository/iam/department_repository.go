package iam

import (
	"context"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/internal/repository"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DepartmentRepository interface {
	FindByID(ctx context.Context, id uint64) (*iammodel.AdminDepartment, error)
	FindByName(ctx context.Context, name string) (*iammodel.AdminDepartment, error)
	ListAll(ctx context.Context) ([]iammodel.AdminDepartment, error)
	ListChildren(ctx context.Context, parentID uint64) ([]iammodel.AdminDepartment, error)
	Create(ctx context.Context, dept *iammodel.AdminDepartment) error
	Update(ctx context.Context, dept *iammodel.AdminDepartment) error
	DeleteByID(ctx context.Context, id uint64) error
}

type departmentRepository struct {
	model iammodel.AdminDepartmentModel
	conn  sqlx.SqlConn
}

func NewDepartmentRepository(repo *repository.Repository) DepartmentRepository {
	return &departmentRepository{model: repo.AdminDepartmentModel, conn: repo.DB}
}

func (r *departmentRepository) FindByID(ctx context.Context, id uint64) (*iammodel.AdminDepartment, error) {
	return r.model.FindOne(ctx, id)
}

// FindByName 按名称查部门。admin_department 没有 name 唯一键（历史遗留，见
// docs/changelog/archive-backend.md 第 23 节同类问题记录），同名重复只取 id 最小的一条。
func (r *departmentRepository) FindByName(ctx context.Context, name string) (*iammodel.AdminDepartment, error) {
	var list []iammodel.AdminDepartment
	sqlStr, args, err := sq.Select("*").From("`admin_department`").
		Where(sq.Eq{"deleted_at": 0, "name": name}).
		OrderBy("id ASC").Limit(1).ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowsCtx(ctx, &list, sqlStr, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql执行有误", err)
	}
	if len(list) == 0 {
		return nil, iammodel.ErrNotFound
	}
	return &list[0], nil
}

func (r *departmentRepository) ListAll(ctx context.Context) ([]iammodel.AdminDepartment, error) {
	list, _, err := r.model.FindPage(ctx, 1, 10000)
	return list, err
}

func (r *departmentRepository) ListChildren(ctx context.Context, parentID uint64) ([]iammodel.AdminDepartment, error) {
	var list []iammodel.AdminDepartment
	sqlStr, args, err := sq.Select("*").From("`admin_department`").
		Where(sq.Eq{"deleted_at": 0, "parent_id": parentID}).
		OrderBy("order_num ASC", "id ASC").ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowsCtx(ctx, &list, sqlStr, args...); err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql执行有误", err)
	}
	return list, nil
}

func (r *departmentRepository) Create(ctx context.Context, dept *iammodel.AdminDepartment) error {
	result, err := r.model.Insert(ctx, dept)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	dept.Id = uint64(id)
	return nil
}

func (r *departmentRepository) Update(ctx context.Context, dept *iammodel.AdminDepartment) error {
	return r.model.Update(ctx, dept)
}

func (r *departmentRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}
