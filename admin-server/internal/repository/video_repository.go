package repository

import (
	"context"
	"time"

	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type VideoRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.Video, error)
	FindByUuid(ctx context.Context, uuid string) (*model.Video, error)
	FindPage(ctx context.Context, page, pageSize int64, keyword string, sourceType int64) ([]model.Video, int64, error)
	FindAllByType(ctx context.Context, sourceType int64) ([]model.Video, error)
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

func (r *videoRepository) FindByUuid(ctx context.Context, uuid string) (*model.Video, error) {
	sql, args, err := sq.Select("id", "uuid", "name", "cover", "god_num", "duration", "play_url", "xlzz_urls", "description", "type", "created_at", "updated_at", "deleted_at").
		From("`video`").
		Where(sq.And{
			sq.Eq{"uuid": uuid},
			sq.Eq{"deleted_at": 0},
		}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var video model.Video
	err = r.conn.QueryRowCtx(ctx, &video, sql, args...)
	if err != nil {
		return nil, err
	}
	return &video, nil
}

func (r *videoRepository) FindPage(ctx context.Context, page, pageSize int64, keyword string, sourceType int64) ([]model.Video, int64, error) {
	// 如果有关键词或类型筛选，需要自定义查询
	if keyword != "" || sourceType > 0 {
		return r.findPageWithFilter(ctx, page, pageSize, keyword, sourceType)
	}
	// 否则使用生成的分页方法
	return r.model.FindPage(ctx, page, pageSize)
}

func (r *videoRepository) findPageWithFilter(ctx context.Context, page, pageSize int64, keyword string, sourceType int64) ([]model.Video, int64, error) {
	conditions := sq.And{sq.Eq{"deleted_at": 0}}

	if keyword != "" {
		keywordPattern := "%" + keyword + "%"
		conditions = append(conditions, sq.Or{
			sq.Like{"name": keywordPattern},
			sq.Like{"description": keywordPattern},
		})
	}

	if sourceType > 0 {
		conditions = append(conditions, sq.Eq{"type": sourceType})
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`video`").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
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

	listSQL, listArgs, err := sq.Select("id", "uuid", "name", "cover", "god_num", "duration", "play_url", "xlzz_urls", "description", "type", "created_at", "updated_at", "deleted_at").
		From("`video`").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []model.Video
	err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *videoRepository) FindAllByType(ctx context.Context, sourceType int64) ([]model.Video, error) {
	conditions := sq.And{
		sq.Eq{"deleted_at": 0},
		sq.Eq{"type": sourceType},
	}

	sql, args, err := sq.Select("id", "uuid", "name", "cover", "god_num", "duration", "play_url", "xlzz_urls", "description", "type", "created_at", "updated_at", "deleted_at").
		From("`video`").
		Where(conditions).
		OrderBy("id DESC").
		ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []model.Video
	err = r.conn.QueryRowsCtx(ctx, &list, sql, args...)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *videoRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *videoRepository) Create(ctx context.Context, video *model.Video) error {
	// 使用 squirrel 手动构建插入语句，确保所有字段都正确插入
	// 因为生成的 Insert 方法可能不包含所有新增字段（uuid, god_num, xlzz_urls, type）
	now := time.Now().Unix()
	if video.CreatedAt == 0 {
		video.CreatedAt = now
	}
	if video.UpdatedAt == 0 {
		video.UpdatedAt = now
	}

	// 构建插入语句
	insert := sq.Insert("`video`").
		Columns("`uuid`", "`name`", "`cover`", "`god_num`", "`duration`", "`play_url`", "`xlzz_urls`", "`description`", "`type`", "`deleted_at`", "`created_at`", "`updated_at`").
		Values(
			video.Uuid,
			video.Name,
			video.Cover,
			video.GodNum,
			video.Duration,
			video.PlayUrl,
			video.XlzzUrls,
			video.Description,
			video.Type,
			video.DeletedAt,
			video.CreatedAt,
			video.UpdatedAt,
		)

	sql, args, err := insert.ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	_, err = r.conn.ExecCtx(ctx, sql, args...)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "插入视频失败", err)
	}

	return nil
}

func (r *videoRepository) Update(ctx context.Context, video *model.Video) error {
	return r.model.Update(ctx, video)
}
