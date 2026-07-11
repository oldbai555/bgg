package monitoring

import (
	"postapocgame/admin-server/internal/repository"
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"postapocgame/admin-server/pkg/errs"
	monitoringmodel "postapocgame/admin-server/internal/model/monitoring"
)

type AuditLogRepository interface {
	FindByID(ctx context.Context, id uint64) (*monitoringmodel.AuditLog, error)
	FindPage(ctx context.Context, page, pageSize int64, userId uint64, username, auditType, auditObject, startTime, endTime string) ([]monitoringmodel.AuditLog, int64, error)
	Create(ctx context.Context, log *monitoringmodel.AuditLog) error
}

type auditLogRepository struct {
	model monitoringmodel.AuditLogModel
	conn  sqlx.SqlConn
}

func NewAuditLogRepository(repo *repository.Repository) AuditLogRepository {
	return &auditLogRepository{model: repo.AuditLogModel, conn: repo.DB}
}

func (r *auditLogRepository) FindByID(ctx context.Context, id uint64) (*monitoringmodel.AuditLog, error) {
	return r.model.FindOne(ctx, id)
}

func (r *auditLogRepository) FindPage(ctx context.Context, page, pageSize int64, userId uint64, username, auditType, auditObject, startTime, endTime string) ([]monitoringmodel.AuditLog, int64, error) {
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

	// 构建查询条件
	conditions := sq.And{sq.Eq{"deleted_at": 0}}

	if userId > 0 {
		conditions = append(conditions, sq.Eq{"user_id": userId})
	}
	if username != "" {
		conditions = append(conditions, sq.Like{"username": "%" + username + "%"})
	}
	if auditType != "" {
		conditions = append(conditions, sq.Eq{"audit_type": auditType})
	}
	if auditObject != "" {
		conditions = append(conditions, sq.Eq{"audit_object": auditObject})
	}
	if startTime != "" {
		// 解析时间字符串为时间戳
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			conditions = append(conditions, sq.GtOrEq{"created_at": t.Unix()})
		}
	}
	if endTime != "" {
		// 解析时间字符串为时间戳
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			conditions = append(conditions, sq.LtOrEq{"created_at": t.Unix()})
		}
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("audit_log").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// 查询列表
	var list []monitoringmodel.AuditLog
	listSQL, listArgs, err := sq.Select("*").
		From("audit_log").
		Where(conditions).
		OrderBy("id DESC").
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

func (r *auditLogRepository) Create(ctx context.Context, log *monitoringmodel.AuditLog) error {
	if log == nil {
		return fmt.Errorf("审计日志数据为空")
	}

	// 设置时间戳
	now := time.Now().Unix()
	if log.CreatedAt == 0 {
		log.CreatedAt = now
	}
	if log.UpdatedAt == 0 {
		log.UpdatedAt = now
	}
	if log.DeletedAt == 0 {
		log.DeletedAt = 0
	}

	result, err := r.model.Insert(ctx, log)
	if err != nil {
		return fmt.Errorf("插入审计日志失败: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取审计日志自增 ID 失败: %w", err)
	}
	log.Id = uint64(id)

	return nil
}
