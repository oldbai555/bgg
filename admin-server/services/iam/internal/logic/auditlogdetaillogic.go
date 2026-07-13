package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iam"
	"postapocgame/admin-server/services/iam/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditLogDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuditLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditLogDetailLogic {
	return &AuditLogDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuditLogDetailLogic) AuditLogDetail(in *iam.AuditLogDetailRequest) (*iam.AuditLogItem, error) {
	if in == nil || in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "审计日志ID不能为空"))
	}

	log, err := l.svcCtx.Domain.Monitoring.AuditLog.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询审计日志详情失败", err))
	}
	if log == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "审计日志不存在"))
	}

	auditDetail := ""
	if log.AuditDetail.Valid {
		auditDetail = log.AuditDetail.String
	}

	return &iam.AuditLogItem{
		Id:          log.Id,
		UserId:      log.UserId,
		Username:    log.Username,
		AuditType:   log.AuditType,
		AuditObject: log.AuditObject,
		AuditDetail: auditDetail,
		IpAddress:   log.IpAddress,
		UserAgent:   log.UserAgent,
		CreatedAt:   log.CreatedAt,
	}, nil
}
