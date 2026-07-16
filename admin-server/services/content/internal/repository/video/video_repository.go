package video

import (
	"postapocgame/admin-server/services/content/internal/repository"
	"context"
	"database/sql"

	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	videomodel "postapocgame/admin-server/services/content/internal/model/video"
)

type VideoRepository interface {
	FindByID(ctx context.Context, id uint64) (*videomodel.Video, error)
	FindByUuid(ctx context.Context, uuid string) (*videomodel.Video, error)
	FindPage(ctx context.Context, page, pageSize int64, keyword string, sourceType int64) ([]videomodel.Video, int64, error)
	FindAllByType(ctx context.Context, sourceType int64) ([]videomodel.Video, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, video *videomodel.Video) error
	Update(ctx context.Context, video *videomodel.Video) error
}

type videoRepository struct {
	model videomodel.VideoModel
	conn  sqlx.SqlConn
}

func NewVideoRepository(store *repository.Store) VideoRepository {
	return &videoRepository{model: store.VideoModel, conn: store.DB}
}

func (r *videoRepository) FindByID(ctx context.Context, id uint64) (*videomodel.Video, error) {
	return r.model.FindOne(ctx, id)
}

func (r *videoRepository) FindByUuid(ctx context.Context, uuid string) (*videomodel.Video, error) {
	return r.model.FindOneByUuid(ctx, sql.NullString{String: uuid, Valid: uuid != ""})
}

func (r *videoRepository) FindPage(ctx context.Context, page, pageSize int64, keyword string, sourceType int64) ([]videomodel.Video, int64, error) {
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

	var list []videomodel.Video
	err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *videoRepository) FindAllByType(ctx context.Context, sourceType int64) ([]videomodel.Video, error) {
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

	var list []videomodel.Video
	err = r.conn.QueryRowsCtx(ctx, &list, sql, args...)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *videoRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *videoRepository) Create(ctx context.Context, video *videomodel.Video) error {
	result, err := r.model.Insert(ctx, video)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "插入视频失败", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "获取视频自增 ID 失败", err)
	}
	video.Id = uint64(id)

	return nil
}

func (r *videoRepository) Update(ctx context.Context, video *videomodel.Video) error {
	return r.model.Update(ctx, video)
}
