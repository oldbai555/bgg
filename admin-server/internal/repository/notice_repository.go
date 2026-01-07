package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"
)

type NoticeRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.AdminNotice, error)
	FindPage(ctx context.Context, page, pageSize int64, title string, noticeType, status int64) ([]model.AdminNotice, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, notice *model.AdminNotice) error
	Update(ctx context.Context, notice *model.AdminNotice) error
	// FindPublishedNotReadByUser 查找已发布且用户未读的公告（通过检查通知表）
	FindPublishedNotReadByUser(ctx context.Context, userID uint64) ([]model.AdminNotice, error)
}

type noticeRepository struct {
	model model.AdminNoticeModel
	conn  sqlx.SqlConn
}

func NewNoticeRepository(repo *Repository) NoticeRepository {
	return &noticeRepository{model: repo.AdminNoticeModel, conn: repo.DB}
}

func (r *noticeRepository) FindByID(ctx context.Context, id uint64) (*model.AdminNotice, error) {
	return r.model.FindOne(ctx, id)
}

func (r *noticeRepository) FindPage(ctx context.Context, page, pageSize int64, title string, noticeType, status int64) ([]model.AdminNotice, int64, error) {
	// 构建查询条件
	conditions := sq.And{sq.Eq{"deleted_at": 0}}

	if title != "" {
		conditions = append(conditions, sq.Like{"title": "%" + title + "%"})
	}
	if noticeType > 0 {
		conditions = append(conditions, sq.Eq{"type": noticeType})
	}
	// status: -1表示未传入（不筛选），1表示草稿，2表示已发布
	if status > 0 {
		conditions = append(conditions, sq.Eq{"status": status})
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`admin_notice`").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	var list []model.AdminNotice
	offset := (page - 1) * pageSize
	listSQL, listArgs, err := sq.Select("id", "title", "content", "type", "status", "publish_time", "created_by", "created_at", "updated_at", "deleted_at").
		From("`admin_notice`").
		Where(conditions).
		OrderBy("publish_time DESC", "created_at DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...)
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *noticeRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *noticeRepository) Create(ctx context.Context, notice *model.AdminNotice) error {
	_, err := r.model.Insert(ctx, notice)
	return err
}

func (r *noticeRepository) Update(ctx context.Context, notice *model.AdminNotice) error {
	return r.model.Update(ctx, notice)
}

func (r *noticeRepository) FindPublishedNotReadByUser(ctx context.Context, userID uint64) ([]model.AdminNotice, error) {
	// 查询已发布且用户未读的公告
	// 通过 LEFT JOIN 查找已发布但用户没有未读通知的公告
	// 状态：1=草稿，2=已发布
	query := `
		SELECT n.id, n.title, n.content, n.type, n.status, n.publish_time, n.created_by, n.created_at, n.updated_at, n.deleted_at
		FROM admin_notice n
		LEFT JOIN admin_notification notif ON n.id = notif.source_id 
			AND notif.source_type = 'notice' 
			AND notif.user_id = ? 
			AND notif.deleted_at = 0
		WHERE n.status = 2 
			AND n.deleted_at = 0
			AND (notif.id IS NULL OR notif.read_status = 1)
		ORDER BY n.publish_time DESC, n.created_at DESC
	`
	var list []model.AdminNotice
	err := r.conn.QueryRowsCtx(ctx, &list, query, userID)
	if err != nil {
		return nil, err
	}
	return list, nil
}
