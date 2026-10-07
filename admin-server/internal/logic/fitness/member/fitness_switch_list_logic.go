// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package member

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

type FitnessSwitchListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessSwitchListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessSwitchListLogic {
	return &FitnessSwitchListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessSwitchListLogic) FitnessSwitchList(req *types.FitnessSwitchListReq) (resp *types.FitnessSwitchListResp, err error) {
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)
	rpcResp, err := l.svcCtx.IamRPC.FitnessSwitchList(l.ctx, &iamclient.FitnessSwitchListRequest{UserId: req.UserId, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, errs.WrapGRPCError("查询切换记录失败", err)
	}
	return &types.FitnessSwitchListResp{Total: rpcResp.Total, List: fitnessconv.SwitchItems(rpcResp.List)}, nil
}
