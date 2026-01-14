// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package metric_admin

import (
	"context"
	"database/sql"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetricStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMetricStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetricStatsLogic {
	return &MetricStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MetricStatsLogic) MetricStats(req *types.MetricStatsReq) (resp *types.MetricStatsResp, err error) {
	// 参数校验和默认值处理
	module := strings.TrimSpace(req.Module)
	if module == "" {
		return nil, errs.New(errs.CodeBadRequest, "module不能为空")
	}

	bizID := req.BizId

	// 日期处理：默认今天，格式 YYYYMMDD
	day := strings.TrimSpace(req.Day)
	if day == "" {
		day = time.Now().Format("20060102")
	} else {
		// 验证日期格式（YYYYMMDD）
		if len(day) != 8 {
			return nil, errs.New(errs.CodeBadRequest, "日期格式错误，应为YYYYMMDD")
		}
		_, err := time.Parse("20060102", day)
		if err != nil {
			return nil, errs.New(errs.CodeBadRequest, "日期格式错误，应为YYYYMMDD")
		}
	}

	// 从 MySQL 的 metric_daily_stats 表中读取聚合结果
	// 约定：
	// - bizId > 0 时：按 module + bizId + day 精确查询
	// - bizId == 0 时：仅按 module + day 维度做整体汇总（SUM 聚合），忽略 bizId 条件
	conditions := sq.And{
		sq.Eq{"module": module},
		sq.Eq{"day": day},
		sq.Eq{"deleted_at": 0},
	}

	var builder sq.SelectBuilder
	if bizID > 0 {
		conditions = append(conditions, sq.Eq{"biz_id": bizID})
		builder = sq.Select("`module`", "`biz_id`", "`day`", "`pv`", "`uv`", "`vv`", "`ip`").
			From("`metric_daily_stats`").
			Where(conditions).
			Limit(1)
	} else {
		// bizId == 0：按 module + day 汇总，返回汇总后的 pv/uv/vv/ip，biz_id 固定为 0
		builder = sq.Select(
			"`module`",
			"0 AS `biz_id`",
			"`day`",
			"SUM(`pv`) AS `pv`",
			"SUM(`uv`) AS `uv`",
			"SUM(`vv`) AS `vv`",
			"SUM(`ip`) AS `ip`",
		).
			From("`metric_daily_stats`").
			Where(conditions).
			GroupBy("`module`", "`day`").
			Limit(1)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "metric_daily_stats 查询sql生成失败", err)
	}

	var row struct {
		Module string `db:"module"`
		BizID  uint64 `db:"biz_id"`
		Day    string `db:"day"`
		Pv     int64  `db:"pv"`
		Uv     int64  `db:"uv"`
		Vv     int64  `db:"vv"`
		Ip     int64  `db:"ip"`
	}

	err = l.svcCtx.Repository.DB.QueryRowCtx(l.ctx, &row, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			// 没有记录时返回 0 值
			return &types.MetricStatsResp{
				Module: module,
				BizId:  bizID,
				Day:    day,
				Pv:     0,
				Uv:     0,
				Vv:     0,
				Ip:     0,
			}, nil
		}
		return nil, errs.Wrap(errs.CodeBadDB, "metric_daily_stats 查询失败", err)
	}

	return &types.MetricStatsResp{
		Module: row.Module,
		BizId:  row.BizID,
		Day:    row.Day,
		Pv:     row.Pv,
		Uv:     row.Uv,
		Vv:     row.Vv,
		Ip:     row.Ip,
	}, nil
}
