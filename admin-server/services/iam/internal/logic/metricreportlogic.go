package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/consts"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetricReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMetricReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetricReportLogic {
	return &MetricReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Metric / Monitor
func (l *MetricReportLogic) MetricReport(in *iam.MetricReportRequest) (*iam.Empty, error) {
	module := strings.TrimSpace(in.Module)
	if module == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "module不能为空"))
	}
	if !isValidMetricModule(module) {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "module不合法"))
	}

	rdb := l.svcCtx.Repository.Redis
	if rdb == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeInternalError, "redis未初始化"))
	}

	now := time.Now()
	day := now.Format("20060102")
	bizID := in.BizId

	baseKey := fmt.Sprintf("metric:%s:%d:%s", module, bizID, day)

	clientIP := in.ClientIp
	userAgent := in.UserAgent

	var deltaPv int64 = 1
	var deltaVv int64 = 1
	var deltaUv int64
	var deltaIp int64

	_, _ = rdb.IncrCtx(l.ctx, baseKey+":pv")
	_ = rdb.ExpireCtx(l.ctx, baseKey+":pv", int((8 * 24 * time.Hour).Seconds()))

	_, _ = rdb.IncrCtx(l.ctx, baseKey+":vv")
	_ = rdb.ExpireCtx(l.ctx, baseKey+":vv", int((8 * 24 * time.Hour).Seconds()))

	if clientIP != "" {
		visitorKey := fmt.Sprintf("%s|%s", clientIP, userAgent)
		visitorHash := fmt.Sprintf("%x", md5.Sum([]byte(visitorKey)))
		uvKey := baseKey + ":uv:" + visitorHash

		exists, err := rdb.Exists(uvKey)
		if err == nil && !exists {
			_ = rdb.Setex(uvKey, "1", int((8 * 24 * time.Hour).Seconds()))
			_, _ = rdb.IncrCtx(l.ctx, baseKey+":uv")
			_ = rdb.ExpireCtx(l.ctx, baseKey+":uv", int((8 * 24 * time.Hour).Seconds()))
			deltaUv = 1
		}
	}

	if clientIP != "" {
		ipKey := baseKey + ":ip:" + clientIP

		exists, err := rdb.Exists(ipKey)
		if err == nil && !exists {
			_ = rdb.Setex(ipKey, "1", int((8 * 24 * time.Hour).Seconds()))
			_, _ = rdb.IncrCtx(l.ctx, baseKey+":ip")
			_ = rdb.ExpireCtx(l.ctx, baseKey+":ip", int((8 * 24 * time.Hour).Seconds()))
			deltaIp = 1
		}
	}

	go func(module string, bizID uint64, day string, dpv, duv, dvv, dip int64) {
		ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
		defer cancel()

		if err := l.svcCtx.Domain.Monitoring.Metric.UpsertDailyStats(ctx, module, bizID, day, dpv, duv, dvv, dip); err != nil {
			logx.Errorf("metric_daily_stats upsert失败: module=%s bizId=%d day=%s err=%v", module, bizID, day, err)
		}
	}(module, bizID, day, deltaPv, deltaUv, deltaVv, deltaIp)

	return &iam.Empty{}, nil
}

func isValidMetricModule(module string) bool {
	switch module {
	case consts.MetricModuleBlogArticleList,
		consts.MetricModuleBlogArticleDetail,
		consts.MetricModuleVideoList,
		consts.MetricModuleVideoDetail:
		return true
	default:
		return false
	}
}
