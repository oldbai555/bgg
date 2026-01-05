package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model"
)

type VideoRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.Video, error)
	FindPage(ctx context.Context, page, pageSize int64, keyword string) ([]model.Video, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, video *model.Video) error
	Update(ctx context.Context, video *model.Video) error
}

type videoRepository struct {
	model model.VideoModel
	conn  sqlx.SqlConn
}

func NewVideoRepository(repo *Repository) VideoRepository {
	return &videoRepository{model: repo.VideoModel, conn: repo.DB}
}

func (r *videoRepository) FindByID(ctx context.Context, id uint64) (*model.Video, error) {
	return r.model.FindOne(ctx, id)
}

func (r *videoRepository) FindPage(ctx context.Context, page, pageSize int64, keyword string) ([]model.Video, int64, error) {
	// 如果有关键词筛选，需要自定义查询
	if keyword != "" {
		return r.findPageWithFilter(ctx, page, pageSize, keyword)
	}
	// 否则使用生成的分页方法
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *videoRepository) findPageWithFilter(ctx context.Context, page, pageSize int64, keyword string) ([]model.Video, int64, error) {
	var whereConditions []string
	var args []interface{}

	whereConditions = append(whereConditions, "deleted_at = 0")

	if keyword != "" {
		whereConditions = append(whereConditions, "(name LIKE ? OR description LIKE ?)")
		keywordPattern := "%" + keyword + "%"
		args = append(args, keywordPattern, keywordPattern)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// 查询总数
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM `video` WHERE %s", whereClause)
	err := r.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	query := fmt.Sprintf("SELECT id, name, cover, duration, play_url, description, created_at, updated_at, deleted_at FROM `video` WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, pageSize, offset)

	var list []model.Video
	err = r.conn.QueryRowsCtx(ctx, &list, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *videoRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *videoRepository) Create(ctx context.Context, video *model.Video) error {
	_, err := r.model.Insert(ctx, video)
	return err
}

func (r *videoRepository) Update(ctx context.Context, video *model.Video) error {
	return r.model.Update(ctx, video)
}
