package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuditLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditLogListLogic {
	return &AuditLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuditLogListLogic) AuditLogList(in *iam.AuditLogListRequest) (*iam.AuditLogListResponse, error) {
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

	list, total, err := l.svcCtx.Domain.Monitoring.AuditLog.FindPage(
		l.ctx, page, pageSize, in.UserId, in.Username, in.AuditType, in.AuditObject, in.StartTime, in.EndTime,
	)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询审计日志列表失败", err))
	}

	items := make([]*iam.AuditLogItem, 0, len(list))
	for _, log := range list {
		auditDetail := ""
		if log.AuditDetail.Valid {
			auditDetail = log.AuditDetail.String
		}
		items = append(items, &iam.AuditLogItem{
			Id:          log.Id,
			UserId:      log.UserId,
			Username:    log.Username,
			AuditType:   log.AuditType,
			AuditObject: log.AuditObject,
			AuditDetail: auditDetail,
			IpAddress:   log.IpAddress,
			UserAgent:   log.UserAgent,
			CreatedAt:   log.CreatedAt,
		})
	}

	return &iam.AuditLogListResponse{List: items, Total: total}, nil
}
