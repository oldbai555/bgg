package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"postapocgame/admin-server/pkg/errs"
)

// MetricRepository 负责 metric_daily_stats 日统计表的数据访问。
type MetricRepository interface {
	// UpsertDailyStats 按 module/bizId/day 聚合累加 PV/UV/VV/IP。
	// 传入的是「增量」，内部使用 INSERT ... ON DUPLICATE KEY UPDATE 做自增。
	UpsertDailyStats(ctx context.Context, module string, bizID uint64, day string, deltaPv, deltaUv, deltaVv, deltaIp int64) error
}

type metricRepository struct {
	repo *Repository
}

func NewMetricRepository(repo *Repository) MetricRepository {
	return &metricRepository{repo: repo}
}

func (r *metricRepository) UpsertDailyStats(
	ctx context.Context,
	module string,
	bizID uint64,
	day string,
	deltaPv, deltaUv, deltaVv, deltaIp int64,
) error {
	// 使用 squirrel 构建 INSERT ... ON DUPLICATE KEY UPDATE
	builder := sq.Insert("`metric_daily_stats`").
		Columns("`module`", "`biz_id`", "`day`", "`pv`", "`uv`", "`vv`", "`ip`", "`created_at`", "`updated_at`", "`deleted_at`").
		Values(module, bizID, day, max64(deltaPv, 0), max64(deltaUv, 0), max64(deltaVv, 0), max64(deltaIp, 0), sq.Expr("UNIX_TIMESTAMP()"), sq.Expr("UNIX_TIMESTAMP()"), 0).
		Suffix("ON DUPLICATE KEY UPDATE " +
			"`pv` = `pv` + VALUES(`pv`), " +
			"`uv` = `uv` + VALUES(`uv`), " +
			"`vv` = `vv` + VALUES(`vv`), " +
			"`ip` = `ip` + VALUES(`ip`), " +
			"`updated_at` = UNIX_TIMESTAMP(), " +
			"`deleted_at` = 0")

	sql, args, err := builder.ToSql()
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "metric_daily_stats upsert sql生成有误", err)
	}

	_, err = r.repo.DB.ExecCtx(ctx, sql, args...)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "metric_daily_stats upsert执行失败", err)
	}
	return nil
}

// max64 用于保证增量为非负数，避免意外的负增量破坏统计。
func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
