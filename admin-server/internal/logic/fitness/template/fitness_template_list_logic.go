// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package template

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

type FitnessTemplateListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFitnessTemplateListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FitnessTemplateListLogic {
	return &FitnessTemplateListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FitnessTemplateListLogic) FitnessTemplateList(req *types.FitnessTemplateListReq) (resp *types.FitnessTemplateListResp, err error) {
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)
	rpcResp, err := l.svcCtx.IamRPC.FitnessTemplateList(l.ctx, &iamclient.FitnessTemplateListRequest{Page: req.Page, PageSize: req.PageSize, Keyword: req.Keyword})
	if err != nil {
		return nil, errs.WrapGRPCError("查询计划模板失败", err)
	}
	items := make([]types.FitnessTemplateItem, 0, len(rpcResp.List))
	for _, t := range rpcResp.List {
		items = append(items, fitnessconv.TemplateItem(t))
	}
	return &types.FitnessTemplateListResp{Total: rpcResp.Total, List: items}, nil
}
