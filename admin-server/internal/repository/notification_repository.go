package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/pkg/errs"
)

type NotificationRepository interface {
	FindByID(ctx context.Context, id uint64) (*model.AdminNotification, error)
	FindPage(ctx context.Context, page, pageSize int64, userID uint64, sourceType string, readStatus int64) ([]model.AdminNotification, int64, error)
	DeleteByID(ctx context.Context, id uint64) error
	Create(ctx context.Context, notification *model.AdminNotification) error
	Update(ctx context.Context, notification *model.AdminNotification) error
	// 全部已读：更新用户的所有未读消息为已读
	MarkAllAsRead(ctx context.Context, userID uint64) error
	// 清除已读：删除用户的所有已读消息
	ClearRead(ctx context.Context, userID uint64) error
}

type notificationRepository struct {
	model model.AdminNotificationModel
	conn  sqlx.SqlConn
}

func NewNotificationRepository(repo *Repository) NotificationRepository {
	return &notificationRepository{model: repo.AdminNotificationModel, conn: repo.DB}
}

func (r *notificationRepository) FindByID(ctx context.Context, id uint64) (*model.AdminNotification, error) {
	return r.model.FindOne(ctx, id)
}

func (r *notificationRepository) FindPage(ctx context.Context, page, pageSize int64, userID uint64, sourceType string, readStatus int64) ([]model.AdminNotification, int64, error) {
	// 构建查询条件
	where := sq.Eq{"deleted_at": 0}

	if userID > 0 {
		where["user_id"] = userID
	}
	if sourceType != "" {
		where["source_type"] = sourceType
	}
	// readStatus <= 0 表示不筛选，readStatus > 0 时添加筛选条件
	// 枚举（字典 read_status）：1 = 未读，2 = 已读
	if readStatus > 0 {
		where["read_status"] = readStatus
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`admin_notification`").Where(where).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	var list []model.AdminNotification
	offset := (page - 1) * pageSize
	listSQL, listArgs, err := sq.Select("id", "user_id", "source_type", "source_id", "title", "content", "read_status", "read_at", "created_at", "updated_at", "deleted_at").
		From("`admin_notification`").
		Where(where).
		OrderBy("created_at DESC").
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

func (r *notificationRepository) DeleteByID(ctx context.Context, id uint64) error {
	return r.model.Delete(ctx, id)
}

func (r *notificationRepository) Create(ctx context.Context, notification *model.AdminNotification) error {
	_, err := r.model.Insert(ctx, notification)
	return err
}

func (r *notificationRepository) Update(ctx context.Context, notification *model.AdminNotification) error {
	return r.model.Update(ctx, notification)
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint64) error {
	now := time.Now().Unix()
	sql, args, err := sq.Update("`admin_notification`").
		Set("`read_status`", 1).
		Set("`read_at`", now).
		Set("`updated_at`", now).
		Where(sq.Eq{"user_id": userID, "read_status": 0, "deleted_at": 0}).
		ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	_, err = r.conn.ExecCtx(ctx, sql, args...)
	return err
}

func (r *notificationRepository) ClearRead(ctx context.Context, userID uint64) error {
	// 软删除所有已读消息
	now := time.Now().Unix()
	sql, args, err := sq.Update("`admin_notification`").
		Set("`deleted_at`", now).
		Set("`updated_at`", now).
		Where(sq.Eq{"user_id": userID, "read_status": 1, "deleted_at": 0}).
		ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	_, err = r.conn.ExecCtx(ctx, sql, args...)
	return err
}
