package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogListLogic {
	return &OperationLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// OperationLog / LoginLog / PerformanceLog / AuditLog（列表/详情/统计；导出见
// services/task 的 GenericExportExecutor，走 pkg/taskcallback，不在本服务）
func (l *OperationLogListLogic) OperationLogList(in *iam.OperationLogListRequest) (*iam.OperationLogListResponse, error) {
	if in == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "请求参数不能为空"))
	}

	page, pageSize := in.Page, in.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	} else if pageSize > 100 {
		pageSize = 100
	}

	list, total, err := l.svcCtx.Domain.Monitoring.OperationLog.FindPage(
		l.ctx, page, pageSize, in.UserId, in.Username, in.OperationType, in.OperationObject, in.Method, in.StartTime, in.EndTime,
	)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询操作日志列表失败", err))
	}

	items := make([]*iam.OperationLogItem, 0, len(list))
	for _, log := range list {
		requestParams := ""
		if log.RequestParams.Valid {
			requestParams = log.RequestParams.String
		}
		items = append(items, &iam.OperationLogItem{
			Id:              log.Id,
			UserId:          log.UserId,
			Username:        log.Username,
			OperationType:   log.OperationType,
			OperationObject: log.OperationObject,
			Method:          log.Method,
			Path:            log.Path,
			RequestParams:   requestParams,
			ResponseCode:    int64(log.ResponseCode),
			ResponseMsg:     log.ResponseMsg,
			IpAddress:       log.IpAddress,
			UserAgent:       log.UserAgent,
			Duration:        int64(log.Duration),
			CreatedAt:       log.CreatedAt,
		})
	}

	return &iam.OperationLogListResponse{Total: total, List: items}, nil
}
