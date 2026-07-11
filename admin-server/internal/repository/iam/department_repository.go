package iam

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	iammodel "postapocgame/admin-server/internal/model/iam"
)

type DepartmentRepository interface {
	FindByID(ctx context.Context, id uint64) (*iammodel.AdminDepartment, error)
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

func (r *departmentRepository) ListAll(ctx context.Context) ([]iammodel.AdminDepartment, error) {
	list, _, err := r.model.FindPage(ctx, 1, 10000)
	return list, err
}

func (r *departmentRepository) ListChildren(ctx context.Context, parentID uint64) ([]iammodel.AdminDepartment, error) {
	var list []iammodel.AdminDepartment
	query := "select * from admin_department where deleted_at = 0 and parent_id = ? order by order_num asc, id asc"
	err := r.conn.QueryRowsCtx(ctx, &list, query, parentID)
	return list, err
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
