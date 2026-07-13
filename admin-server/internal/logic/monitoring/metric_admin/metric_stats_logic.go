// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package metric_admin

import (
	"context"
	"strings"
	"time"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type MetricStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMetricStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MetricStatsLogic {
	return &MetricStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MetricStatsLogic) MetricStats(req *types.MetricStatsReq) (resp *types.MetricStatsResp, err error) {
	module := strings.TrimSpace(req.Module)
	if module == "" {
		return nil, errs.New(errs.CodeBadRequest, "module不能为空")
	}

	day := strings.TrimSpace(req.Day)
	if day == "" {
		day = time.Now().Format("20060102")
	} else {
		if len(day) != 8 {
			return nil, errs.New(errs.CodeBadRequest, "日期格式错误，应为YYYYMMDD")
		}
		if _, err := time.Parse("20060102", day); err != nil {
			return nil, errs.New(errs.CodeBadRequest, "日期格式错误，应为YYYYMMDD")
		}
	}

	rpcResp, err := l.svcCtx.IamRPC.MetricStats(l.ctx, &iamclient.MetricStatsRequest{
		Module: module,
		BizId:  req.BizId,
		Day:    day,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询统计数据失败", err)
	}

	return &types.MetricStatsResp{
		Module: rpcResp.Module,
		BizId:  rpcResp.BizId,
		Day:    rpcResp.Day,
		Pv:     rpcResp.Pv,
		Uv:     rpcResp.Uv,
		Vv:     rpcResp.Vv,
		Ip:     rpcResp.Ip,
	}, nil
}
