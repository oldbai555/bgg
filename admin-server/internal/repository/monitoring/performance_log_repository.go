package monitoring

import (
	"postapocgame/admin-server/internal/repository"
	"context"
	"time"

	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	monitoringmodel "postapocgame/admin-server/internal/model/monitoring"
)

// PerformanceLogRepository 性能监控日志仓库
type PerformanceLogRepository interface {
	// FindPage 分页查询性能日志（预留给后续列表接口使用）
	FindPage(ctx context.Context, page, pageSize int64, method, path string, isSlow int64, statusCode int64, startTime, endTime string) ([]monitoringmodel.AdminPerformanceLog, int64, error)
	// Create 创建一条性能日志记录
	Create(ctx context.Context, log *monitoringmodel.AdminPerformanceLog) error
}

type performanceLogRepository struct {
	model monitoringmodel.AdminPerformanceLogModel
	conn  sqlx.SqlConn
}

// NewPerformanceLogRepository 创建性能日志仓库
func NewPerformanceLogRepository(repo *repository.Repository) PerformanceLogRepository {
	return &performanceLogRepository{
		model: repo.AdminPerformanceLogModel,
		conn:  repo.DB,
	}
}

// FindPage 分页查询性能日志
func (r *performanceLogRepository) FindPage(ctx context.Context, page, pageSize int64, method, path string, isSlow int64, statusCode int64, startTime, endTime string) ([]monitoringmodel.AdminPerformanceLog, int64, error) {
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

	// 构建查询条件（使用 squirrel）
	conditions := sq.And{sq.Eq{"deleted_at": 0}}
	if method != "" {
		conditions = append(conditions, sq.Eq{"method": method})
	}
	if path != "" {
		conditions = append(conditions, sq.Like{"path": "%" + path + "%"})
	}
	if isSlow != 0 {
		conditions = append(conditions, sq.Eq{"is_slow": isSlow})
	}
	if statusCode > 0 {
		conditions = append(conditions, sq.Eq{"status_code": statusCode})
	}
	if startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			conditions = append(conditions, sq.GtOrEq{"created_at": t.Unix()})
		}
	}
	if endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			conditions = append(conditions, sq.LtOrEq{"created_at": t.Unix()})
		}
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").
		From("`admin_performance_log`").
		Where(conditions).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	if err := r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...); err != nil {
		return nil, 0, err
	}

	// 查询列表
	listSQL, listArgs, err := sq.Select("*").
		From("`admin_performance_log`").
		Where(conditions).
		OrderBy("id DESC").
		Limit(uint64(pageSize)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}

	var list []monitoringmodel.AdminPerformanceLog
	if err := r.conn.QueryRowsCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// Create 创建一条性能日志记录
func (r *performanceLogRepository) Create(ctx context.Context, log *monitoringmodel.AdminPerformanceLog) error {
	if log == nil {
		return nil
	}
	_, err := r.model.Insert(ctx, log)
	return err
}
