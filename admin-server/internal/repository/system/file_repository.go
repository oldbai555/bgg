package system

import (
	"postapocgame/admin-server/internal/repository"
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	systemmodel "postapocgame/admin-server/internal/model/system"
)

type FileRepository interface {
	FindByID(ctx context.Context, id uint64) (*systemmodel.AdminFile, error)
	FindByName(ctx context.Context, name string) (*systemmodel.AdminFile, error) // 根据 name（MD5）查找文件
	FindPage(ctx context.Context, page, pageSize int64, name string) ([]systemmodel.AdminFile, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, file *systemmodel.AdminFile) error
	Update(ctx context.Context, file *systemmodel.AdminFile) error
}

type fileRepository struct {
	model systemmodel.AdminFileModel
	conn  sqlx.SqlConn
}

func NewFileRepository(repo *repository.Repository) FileRepository {
	return &fileRepository{model: repo.AdminFileModel, conn: repo.DB}
}

func (r *fileRepository) FindByID(ctx context.Context, id uint64) (*systemmodel.AdminFile, error) {
	return r.model.FindOne(ctx, id)
}

func (r *fileRepository) FindByName(ctx context.Context, name string) (*systemmodel.AdminFile, error) {
	query := "SELECT id, name, original_name, path, base_url, size, mime_type, ext, storage_type, status, created_at, updated_at, deleted_at FROM `admin_file` WHERE `name` = ? AND `deleted_at` = 0 LIMIT 1"
	var result systemmodel.AdminFile
	err := r.conn.QueryRowCtx(ctx, &result, query, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *fileRepository) FindPage(ctx context.Context, page, pageSize int64, name string) ([]systemmodel.AdminFile, int64, error) {
	// 目前生成方法不支持复杂过滤，简单复用生成的分页
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *fileRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *fileRepository) Create(ctx context.Context, file *systemmodel.AdminFile) error {
	result, err := r.model.Insert(ctx, file)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	file.Id = uint64(id)
	return nil
}

func (r *fileRepository) Update(ctx context.Context, file *systemmodel.AdminFile) error {
	return r.model.Update(ctx, file)
}
