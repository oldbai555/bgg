package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOperationLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogDetailLogic {
	return &OperationLogDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OperationLogDetailLogic) OperationLogDetail(in *iam.OperationLogDetailRequest) (*iam.OperationLogItem, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "操作日志ID不能为空"))
	}

	log, err := l.svcCtx.Domain.Monitoring.OperationLog.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询操作日志失败", err))
	}

	requestParams := ""
	if log.RequestParams.Valid {
		requestParams = log.RequestParams.String
	}

	return &iam.OperationLogItem{
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
	}, nil
}
