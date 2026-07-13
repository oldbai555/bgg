// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package metric

import (
	"context"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
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

// MetricReport 薄胶水：解析请求 + 从 context 取客户端 IP/UA -> 调 IamRPC，实际的 Redis
// 计数与 metric_daily_stats 落库已经搬进 services/iam/internal/logic/metricreportlogic.go。
func (l *MetricReportLogic) MetricReport(req *types.MetricReportReq) (resp *types.Response, err error) {
	module := strings.TrimSpace(req.Module)
	if module == "" {
		return nil, errs.New(errs.CodeBadRequest, "module不能为空")
	}

	clientIP := ""
	if ip, ok := l.ctx.Value(ctxKeyClientIP).(string); ok {
		clientIP = ip
	}
	userAgent := ""
	if ua, ok := l.ctx.Value(ctxKeyUserAgent).(string); ok {
		userAgent = ua
	}

	_, err = l.svcCtx.IamRPC.MetricReport(l.ctx, &iamclient.MetricReportRequest{
		Module:    module,
		BizId:     req.BizId,
		ClientIp:  clientIP,
		UserAgent: userAgent,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("上报统计失败", err)
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "ok"}, nil
}
