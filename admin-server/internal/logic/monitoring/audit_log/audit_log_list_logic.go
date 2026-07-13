// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package audit_log

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
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

	page, pageSize := logicutil.NormalizePage(int64(req.Page), int64(req.PageSize), 20, 100)

	rpcResp, err := l.svcCtx.IamRPC.AuditLogList(l.ctx, &iamclient.AuditLogListRequest{
		Page:        page,
		PageSize:    pageSize,
		UserId:      req.UserId,
		Username:    req.Username,
		AuditType:   req.AuditType,
		AuditObject: req.AuditObject,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询审计日志列表失败", err)
	}

	items := make([]types.AuditLogItem, 0, len(rpcResp.List))
	for _, log := range rpcResp.List {
		items = append(items, types.AuditLogItem{
			Id:          log.Id,
			UserId:      log.UserId,
			Username:    log.Username,
			AuditType:   log.AuditType,
			AuditObject: log.AuditObject,
			AuditDetail: log.AuditDetail,
			IpAddress:   log.IpAddress,
			UserAgent:   log.UserAgent,
			CreatedAt:   log.CreatedAt,
		})
	}

	return &types.AuditLogListResp{
		List:     items,
		Total:    rpcResp.Total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
