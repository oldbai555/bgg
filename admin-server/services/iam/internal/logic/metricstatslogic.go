package logic

import (
	"context"
	"database/sql"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetricStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMetricStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetricStatsLogic {
	return &MetricStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MetricStatsLogic) MetricStats(in *iam.MetricStatsRequest) (*iam.MetricStatsResponse, error) {
	module := strings.TrimSpace(in.Module)
	if module == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "module不能为空"))
	}

	bizID := in.BizId

	day := strings.TrimSpace(in.Day)
	if day == "" {
		day = time.Now().Format("20060102")
	} else {
		if len(day) != 8 {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "日期格式错误，应为YYYYMMDD"))
		}
		if _, err := time.Parse("20060102", day); err != nil {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "日期格式错误，应为YYYYMMDD"))
		}
	}

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
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "metric_daily_stats 查询sql生成失败", err))
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
			return &iam.MetricStatsResponse{Module: module, BizId: bizID, Day: day}, nil
		}
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "metric_daily_stats 查询失败", err))
	}

	return &iam.MetricStatsResponse{
		Module: row.Module,
		BizId:  row.BizID,
		Day:    row.Day,
		Pv:     row.Pv,
		Uv:     row.Uv,
		Vv:     row.Vv,
		Ip:     row.Ip,
	}, nil
}
