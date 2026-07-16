package logic

import (
	"context"
	"time"

	"postapocgame/admin-server/services/iam/iam"
	monitoringmodel "postapocgame/admin-server/services/iam/internal/model/monitoring"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecordPerformanceLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecordPerformanceLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecordPerformanceLogLogic {
	return &RecordPerformanceLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RecordPerformanceLog 供 gateway PerformanceMiddleware（仅慢接口/错误请求才调用一次，
// 原语义在 internal/middleware/performancemiddleware.go 里是"异步写入"，gateway 侧继续
// 用 go func() 包一层这次调用，这里只做同步落库）。
func (l *RecordPerformanceLogLogic) RecordPerformanceLog(in *iam.RecordPerformanceLogRequest) (*iam.Empty, error) {
	now := time.Now().Unix()
	log := &monitoringmodel.AdminPerformanceLog{
		UserId:        in.UserId,
		Username:      in.Username,
		Method:        in.Method,
		Path:          in.Path,
		StatusCode:    in.StatusCode,
		Duration:      in.Duration,
		IsSlow:        in.IsSlow,
		SlowThreshold: in.SlowThreshold,
		IpAddress:     in.IpAddress,
		UserAgent:     in.UserAgent,
		ErrorMsg:      in.ErrorMsg,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     0,
	}

	if err := l.svcCtx.Domain.Monitoring.PerformanceLog.Create(l.ctx, log); err != nil {
		l.Errorf("写入性能监控日志失败: method=%s, path=%s, error=%v", in.Method, in.Path, err)
	}

	return &iam.Empty{}, nil
}
