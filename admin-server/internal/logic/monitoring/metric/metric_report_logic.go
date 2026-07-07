// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package metric

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	monitoringrepo "postapocgame/admin-server/internal/repository/monitoring"
)

type MetricReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMetricReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetricReportLogic {
	return &MetricReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const (
	// IMPORTANT: keep consistent with handler context keys (plain strings).
	ctxKeyClientIP  = "client_ip"
	ctxKeyUserAgent = "user_agent"
)

func (l *MetricReportLogic) MetricReport(req *types.MetricReportReq) (resp *types.Response, err error) {
	module := strings.TrimSpace(req.Module)
	if module == "" {
		return nil, errs.New(errs.CodeBadRequest, "module不能为空")
	}
	if !isValidMetricModule(module) {
		return nil, errs.New(errs.CodeBadRequest, "module不合法")
	}

	rdb := l.svcCtx.Repository.Redis
	if rdb == nil {
		return nil, errs.New(errs.CodeInternalError, "redis未初始化")
	}

	now := time.Now()
	day := now.Format("20060102")
	bizID := req.BizId

	baseKey := fmt.Sprintf("metric:%s:%d:%s", module, bizID, day)

	// 获取客户端 IP 和 User-Agent
	clientIP := ""
	if ip, ok := l.ctx.Value(ctxKeyClientIP).(string); ok {
		clientIP = ip
	}
	userAgent := ""
	if ua, ok := l.ctx.Value(ctxKeyUserAgent).(string); ok {
		userAgent = ua
	}

	// PV/VV 增量始终为 1
	var deltaPv int64 = 1
	var deltaVv int64 = 1
	var deltaUv int64
	var deltaIp int64

	// PV：累计访问次数（每次请求都增加）
	_, _ = rdb.IncrCtx(l.ctx, baseKey+":pv")
	_ = rdb.ExpireCtx(l.ctx, baseKey+":pv", int((8 * 24 * time.Hour).Seconds()))

	// VV：访问次数（每次请求都增加，与 PV 相同，但保留独立 key 以便未来扩展 session 去重）
	_, _ = rdb.IncrCtx(l.ctx, baseKey+":vv")
	_ = rdb.ExpireCtx(l.ctx, baseKey+":vv", int((8 * 24 * time.Hour).Seconds()))

	// UV：独立访客数（基于 visitorKey = IP + User-Agent 去重）
	if clientIP != "" {
		visitorKey := fmt.Sprintf("%s|%s", clientIP, userAgent)
		visitorHash := fmt.Sprintf("%x", md5.Sum([]byte(visitorKey)))
		uvKey := baseKey + ":uv:" + visitorHash

		// 使用 Exists + Setex 实现去重：如果 key 不存在则设置，存在则忽略
		// 设置过期时间为 8 天，确保跨天统计时不会重复计算
		exists, err := rdb.Exists(uvKey)
		if err == nil && !exists {
			// 首次访问，设置标记并增加 UV 计数
			_ = rdb.Setex(uvKey, "1", int((8 * 24 * time.Hour).Seconds()))
			_, _ = rdb.IncrCtx(l.ctx, baseKey+":uv")
			_ = rdb.ExpireCtx(l.ctx, baseKey+":uv", int((8 * 24 * time.Hour).Seconds()))
			deltaUv = 1
		}
	}

	// IP：独立 IP 数（基于 IP 去重）
	if clientIP != "" {
		ipKey := baseKey + ":ip:" + clientIP

		// 使用 Exists + Setex 实现去重：如果 key 不存在则设置，存在则忽略
		exists, err := rdb.Exists(ipKey)
		if err == nil && !exists {
			// 首次访问，设置标记并增加 IP 计数
			_ = rdb.Setex(ipKey, "1", int((8 * 24 * time.Hour).Seconds()))
			_, _ = rdb.IncrCtx(l.ctx, baseKey+":ip")
			_ = rdb.ExpireCtx(l.ctx, baseKey+":ip", int((8 * 24 * time.Hour).Seconds()))
			deltaIp = 1
		}
	}

	// 将本次增量异步落库到 MySQL 的 metric_daily_stats 表
	go func(module string, bizID uint64, day string, dpv, duv, dvv, dip int64) {
		// 使用独立的 context.Background()，避免 HTTP 请求返回后 context 被取消导致数据库操作失败
		// 设置 800ms 超时，避免数据库慢查询阻塞 goroutine
		ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
		defer cancel()

		repo := monitoringrepo.NewMetricRepository(l.svcCtx.Repository)
		if err := repo.UpsertDailyStats(ctx, module, bizID, day, dpv, duv, dvv, dip); err != nil {
			logx.Errorf("metric_daily_stats upsert失败: module=%s bizId=%d day=%s err=%v", module, bizID, day, err)
		}
	}(module, bizID, day, deltaPv, deltaUv, deltaVv, deltaIp)

	return &types.Response{Code: int(errs.CodeOK), Message: "ok"}, nil
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
