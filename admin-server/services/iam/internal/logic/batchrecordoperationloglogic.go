package logic

import (
	"context"
	"database/sql"

	"postapocgame/admin-server/services/iam/iam"
	monitoringmodel "postapocgame/admin-server/services/iam/internal/model/monitoring"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchRecordOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchRecordOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRecordOperationLogLogic {
	return &BatchRecordOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchRecordOperationLog 供 gateway OperationLogMiddleware/PublicOperationLogMiddleware
// 批量写入操作日志（原逻辑在 admin.go/中间件里直连 Repository，见两个中间件文件的
// writeBatch，逐条 fallback 保留同样的语义）。
func (l *BatchRecordOperationLogLogic) BatchRecordOperationLog(in *iam.BatchRecordOperationLogRequest) (*iam.Empty, error) {
	if len(in.Logs) == 0 {
		return &iam.Empty{}, nil
	}

	logs := make([]*monitoringmodel.AdminOperationLog, 0, len(in.Logs))
	for _, entry := range in.Logs {
		logs = append(logs, &monitoringmodel.AdminOperationLog{
			UserId:          entry.UserId,
			Username:        entry.Username,
			OperationType:   entry.OperationType,
			OperationObject: entry.OperationObject,
			Method:          entry.Method,
			Path:            entry.Path,
			RequestParams:   sql.NullString{String: entry.RequestParams, Valid: entry.RequestParams != ""},
			ResponseCode:    entry.ResponseCode,
			ResponseMsg:     entry.ResponseMsg,
			IpAddress:       entry.IpAddress,
			UserAgent:       entry.UserAgent,
			Duration:        entry.Duration,
			DeletedAt:       0,
		})
	}

	if err := l.svcCtx.Domain.Monitoring.OperationLog.BatchCreate(l.ctx, logs); err != nil {
		l.Errorf("批量写入操作日志失败: count=%d, error: %v", len(logs), err)
		for _, log := range logs {
			if err := l.svcCtx.Domain.Monitoring.OperationLog.Create(l.ctx, log); err != nil {
				l.Errorf("写入操作日志失败: %+v, error: %v", log, err)
			}
		}
	}

	return &iam.Empty{}, nil
}
