// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package metric

import (
	"context"

	"postapocgame/admin-server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
