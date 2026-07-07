package monitoring

import (
	"postapocgame/admin-server/internal/repository"
	"context"
	"fmt"
	"time"

	"postapocgame/admin-server/pkg/errs"

	sq "github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	monitoringmodel "postapocgame/admin-server/internal/model/monitoring"
)

type LoginLogRepository interface {
	FindByID(ctx context.Context, id uint64) (*monitoringmodel.AdminLoginLog, error)
	FindPage(ctx context.Context, page, pageSize int64, userId uint64, username string, status int, startTime, endTime string) ([]monitoringmodel.AdminLoginLog, int64, error)
	Create(ctx context.Context, log *monitoringmodel.AdminLoginLog) error
	// 统计功能
	CountByStatus(ctx context.Context, status int) (int64, error)
	CountToday(ctx context.Context) (int64, error)
	CountTodayByStatus(ctx context.Context, status int) (int64, error)
}

type loginLogRepository struct {
	model monitoringmodel.AdminLoginLogModel
	conn  sqlx.SqlConn
}

func NewLoginLogRepository(repo *repository.Repository) LoginLogRepository {
	return &loginLogRepository{model: repo.AdminLoginLogModel, conn: repo.DB}
}

func (r *loginLogRepository) FindByID(ctx context.Context, id uint64) (*monitoringmodel.AdminLoginLog, error) {
	return r.model.FindOne(ctx, id)
}

func (r *loginLogRepository) FindPage(ctx context.Context, page, pageSize int64, userId uint64, username string, status int, startTime, endTime string) ([]monitoringmodel.AdminLoginLog, int64, error) {
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
	// status <= 0 表示不筛选，status > 0 时添加筛选条件
	// 枚举（字典 login_status）：1 = 成功，2 = 失败
	if status > 0 {
		conditions = append(conditions, sq.Eq{"status": status})
	}
	if startTime != "" {
		// 解析时间字符串为时间戳
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			conditions = append(conditions, sq.GtOrEq{"login_at": t.Unix()})
		}
	}
	if endTime != "" {
		// 解析时间字符串为时间戳
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			conditions = append(conditions, sq.LtOrEq{"login_at": t.Unix()})
		}
	}

	// 查询总数
	var total int64
	countSQL, countArgs, err := sq.Select("COUNT(*)").From("`admin_login_log`").Where(conditions).ToSql()
	if err != nil {
		return nil, 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &total, countSQL, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	var list []monitoringmodel.AdminLoginLog
	listSQL, listArgs, err := sq.Select("*").
		From("`admin_login_log`").
		Where(conditions).
		OrderBy("login_at DESC").
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

func (r *loginLogRepository) Create(ctx context.Context, log *monitoringmodel.AdminLoginLog) error {
	if log == nil {
		return fmt.Errorf("登录日志数据为空")
	}

	// 记录调试信息
	logx.Infof("准备插入登录日志: userId=%d, username=%s, status=%d, message=%s",
		log.UserId, log.Username, log.Status, log.Message)

	result, err := r.model.Insert(ctx, log)
	if err != nil {
		logx.Errorf("插入登录日志失败: userId=%d, username=%s, error: %v",
			log.UserId, log.Username, err)
		return fmt.Errorf("插入登录日志失败: %w", err)
	}

	// 获取插入的 ID（用于调试）
	if result != nil {
		if id, err := result.LastInsertId(); err == nil {
			logx.Infof("成功插入登录日志: id=%d, userId=%d, username=%s", id, log.UserId, log.Username)
		}
	}
	return nil
}

func (r *loginLogRepository) CountByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	var sql string
	var args []interface{}
	var err error
	if status < 0 {
		// status < 0 表示查询所有状态
		sql, args, err = sq.Select("COUNT(*)").From("`admin_login_log`").Where(sq.Eq{"deleted_at": 0}).ToSql()
	} else {
		sql, args, err = sq.Select("COUNT(*)").From("`admin_login_log`").Where(sq.Eq{"deleted_at": 0, "status": status}).ToSql()
	}
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &count, sql, args...)
	return count, err
}

func (r *loginLogRepository) CountToday(ctx context.Context) (int64, error) {
	var count int64
	// 获取今天的开始时间戳（00:00:00）
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	sql, args, err := sq.Select("COUNT(*)").From("`admin_login_log`").Where(sq.And{
		sq.Eq{"deleted_at": 0},
		sq.GtOrEq{"login_at": todayStart},
	}).ToSql()
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &count, sql, args...)
	return count, err
}

func (r *loginLogRepository) CountTodayByStatus(ctx context.Context, status int) (int64, error) {
	var count int64
	// 获取今天的开始时间戳（00:00:00）
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	sql, args, err := sq.Select("COUNT(*)").From("`admin_login_log`").Where(sq.And{
		sq.Eq{"deleted_at": 0, "status": status},
		sq.GtOrEq{"login_at": todayStart},
	}).ToSql()
	if err != nil {
		return 0, errs.Wrap(errs.CodeBadDB, "sql生成有误", err)
	}
	err = r.conn.QueryRowCtx(ctx, &count, sql, args...)
	return count, err
}
