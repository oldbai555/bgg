package logic

import (
	"context"

	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PerformanceLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPerformanceLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PerformanceLogListLogic {
	return &PerformanceLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PerformanceLogListLogic) PerformanceLogList(in *iam.PerformanceLogListRequest) (*iam.PerformanceLogListResponse, error) {
	page, pageSize := in.Page, in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	logs, total, err := l.svcCtx.Domain.Monitoring.PerformanceLog.FindPage(
		l.ctx, page, pageSize, in.Method, in.Path, in.IsSlow, in.StatusCode, in.StartTime, in.EndTime,
	)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	items := make([]*iam.PerformanceLogItem, 0, len(logs))
	for _, lg := range logs {
		items = append(items, &iam.PerformanceLogItem{
			Id:            lg.Id,
			UserId:        lg.UserId,
			Username:      lg.Username,
			Method:        lg.Method,
			Path:          lg.Path,
			StatusCode:    lg.StatusCode,
			Duration:      lg.Duration,
			IsSlow:        lg.IsSlow,
			SlowThreshold: lg.SlowThreshold,
			IpAddress:     lg.IpAddress,
			UserAgent:     lg.UserAgent,
			ErrorMsg:      lg.ErrorMsg,
			CreatedAt:     lg.CreatedAt,
		})
	}

	return &iam.PerformanceLogListResponse{Total: total, List: items}, nil
}
