// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package tip

import (
	"context"

	"postapocgame/admin-server/internal/logic/fitness/fitnessconv"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type FitnessTipListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTipListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTipListLogic {
	return &FitnessTipListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTipListLogic) FitnessTipList(req *types.FitnessTipListReq) (resp *types.FitnessTipListResp, err error) {
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)
	rpcResp, err := l.svcCtx.IamRPC.FitnessTipList(l.ctx, &iamclient.FitnessTipListRequest{Page: req.Page, PageSize: req.PageSize, Keyword: req.Keyword})
	if err != nil {
		return nil, errs.WrapGRPCError("查询通用提示失败", err)
	}
	return &types.FitnessTipListResp{Total: rpcResp.Total, List: fitnessconv.TipItems(rpcResp.List)}, nil
}
