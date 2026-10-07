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

type FitnessMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessMemberListLogic {
	return &FitnessMemberListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessMemberListLogic) FitnessMemberList(req *types.FitnessMemberListReq) (resp *types.FitnessMemberListResp, err error) {
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)
	rpcResp, err := l.svcCtx.IamRPC.FitnessMemberList(l.ctx, &iamclient.FitnessMemberListRequest{Page: req.Page, PageSize: req.PageSize, Keyword: req.Keyword})
	if err != nil {
		return nil, errs.WrapGRPCError("查询使用者列表失败", err)
	}
	items := make([]types.FitnessMemberItem, 0, len(rpcResp.List))
	for _, m := range rpcResp.List {
		items = append(items, fitnessconv.MemberItem(m))
	}
	return &types.FitnessMemberListResp{Total: rpcResp.Total, List: items}, nil
}
