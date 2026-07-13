// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package operation_log

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type OperationLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OperationLogListLogic {
	return &OperationLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OperationLogListLogic) OperationLogList(req *types.OperationLogListReq) (resp *types.OperationLogListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	rpcResp, err := l.svcCtx.IamRPC.OperationLogList(l.ctx, &iamclient.OperationLogListRequest{
		Page:            req.Page,
		PageSize:        req.PageSize,
		UserId:          req.UserId,
		Username:        req.Username,
		OperationType:   req.OperationType,
		OperationObject: req.OperationObject,
		Method:          req.Method,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询操作日志列表失败", err)
	}

	items := make([]types.OperationLogItem, 0, len(rpcResp.List))
	for _, log := range rpcResp.List {
		items = append(items, types.OperationLogItem{
			Id:              log.Id,
			UserId:          log.UserId,
			Username:        log.Username,
			OperationType:   log.OperationType,
			OperationObject: log.OperationObject,
			Method:          log.Method,
			Path:            log.Path,
			RequestParams:   log.RequestParams,
			ResponseCode:    int(log.ResponseCode),
			ResponseMsg:     log.ResponseMsg,
			IpAddress:       log.IpAddress,
			UserAgent:       log.UserAgent,
			Duration:        int(log.Duration),
			CreatedAt:       log.CreatedAt,
		})
	}

	return &types.OperationLogListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
