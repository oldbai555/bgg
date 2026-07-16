// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package demo

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoListLogic {
	return &DemoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoListLogic) DemoList(req *types.DemoListReq) (resp *types.DemoListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.DemoList(l.ctx, &iamclient.DemoListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询演示功能列表失败", err)
	}

	items := make([]types.DemoItem, 0, len(rpcResp.List))
	for _, d := range rpcResp.List {
		items = append(items, types.DemoItem{
			Id:        d.Id,
			Name:      d.Name,
			Status:    d.Status,
			CreatedAt: d.CreatedAt,
		})
	}

	return &types.DemoListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
