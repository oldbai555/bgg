// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package audit_log

import (
	"context"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	monitoringrepo "postapocgame/admin-server/internal/repository/monitoring"
)

type AuditLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditLogListLogic {
	return &AuditLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuditLogListLogic) AuditLogList(req *types.AuditLogListReq) (resp *types.AuditLogListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	// 统一分页参数
	page, pageSize := logicutil.NormalizePage(int64(req.Page), int64(req.PageSize), 20, 100)

	auditLogRepo := monitoringrepo.NewAuditLogRepository(l.svcCtx.Repository)
	list, total, err := auditLogRepo.FindPage(
		l.ctx,
		page,
		pageSize,
		req.UserId,
		req.Username,
		req.AuditType,
		req.AuditObject,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询审计日志列表失败", err)
	}

	items := make([]types.AuditLogItem, 0, len(list))
	for _, log := range list {
		auditDetail := ""
		if log.AuditDetail.Valid {
			auditDetail = log.AuditDetail.String
		}

		items = append(items, types.AuditLogItem{
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

	return &types.AuditLogListResp{
		List:     items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
