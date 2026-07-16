// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package performance_log

import (
	"context"

	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PerformanceLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPerformanceLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PerformanceLogListLogic {
	return &PerformanceLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PerformanceLogListLogic) PerformanceLogList(req *types.PerformanceLogListReq) (resp *types.PerformanceLogListResp, err error) {
	page, pageSize := logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	rpcResp, err := l.svcCtx.IamRPC.PerformanceLogList(l.ctx, &iamclient.PerformanceLogListRequest{
		Page:       page,
		PageSize:   pageSize,
		Method:     req.Method,
		Path:       req.Path,
		IsSlow:     req.IsSlow,
		StatusCode: req.StatusCode,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询性能日志列表失败", err)
	}

	items := make([]types.PerformanceLogItem, 0, len(rpcResp.List))
	for _, lg := range rpcResp.List {
		items = append(items, types.PerformanceLogItem{
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

	return &types.PerformanceLogListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
