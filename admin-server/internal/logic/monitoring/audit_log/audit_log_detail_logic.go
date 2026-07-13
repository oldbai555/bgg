// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package audit_log

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditLogDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditLogDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditLogDetailLogic {
	return &AuditLogDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuditLogDetailLogic) AuditLogDetail(req *types.AuditLogDetailReq) (resp *types.AuditLogDetailResp, err error) {
	if req == nil || req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "审计日志ID不能为空")
	}

	log, err := l.svcCtx.IamRPC.AuditLogDetail(l.ctx, &iamclient.AuditLogDetailRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("查询审计日志详情失败", err)
	}

	return &types.AuditLogDetailResp{
		AuditLogItem: types.AuditLogItem{
			Id:          log.Id,
			UserId:      log.UserId,
			Username:    log.Username,
			AuditType:   log.AuditType,
			AuditObject: log.AuditObject,
			AuditDetail: log.AuditDetail,
			IpAddress:   log.IpAddress,
			UserAgent:   log.UserAgent,
			CreatedAt:   log.CreatedAt,
		},
	}, nil
}
