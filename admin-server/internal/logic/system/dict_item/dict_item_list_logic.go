// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict_item

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictItemListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictItemListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictItemListLogic {
	return &DictItemListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictItemListLogic) DictItemList(req *types.DictItemListReq) (resp *types.DictItemListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.DictItemList(l.ctx, &iamclient.DictItemListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		TypeId:   req.TypeId,
		Label:    req.Label,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询字典项列表失败", err)
	}

	items := make([]types.DictItemItem, 0, len(rpcResp.List))
	for _, di := range rpcResp.List {
		items = append(items, types.DictItemItem{
			Id:        di.Id,
			TypeId:    di.TypeId,
			Label:     di.Label,
			Value:     di.Value,
			Sort:      di.Sort,
			Status:    di.Status,
			Remark:    di.Remark,
			CreatedAt: di.CreatedAt,
		})
	}

	return &types.DictItemListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
