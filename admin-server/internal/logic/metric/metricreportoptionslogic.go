// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package metric

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/svc"
)

type MetricReportOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMetricReportOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetricReportOptionsLogic {
	return &MetricReportOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MetricReportOptionsLogic) MetricReportOptions() error {
	// todo: add your logic here and delete this line

	return nil
}
