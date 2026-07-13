// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictBatchGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictBatchGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictBatchGetLogic {
	return &DictBatchGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictBatchGetLogic) DictBatchGet(req *types.DictBatchGetReq) (resp *types.DictBatchGetResp, err error) {
	if req == nil || len(req.Codes) == 0 {
		return nil, errs.New(errs.CodeBadRequest, "字典类型编码列表不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.DictBatchGet(l.ctx, &iamclient.DictBatchGetRequest{Codes: req.Codes})
	if err != nil {
		return nil, errs.WrapGRPCError("批量查询字典失败", err)
	}

	result := make(map[string]types.DictGetResp, len(rpcResp.Dicts))
	for code, d := range rpcResp.Dicts {
		items := make([]types.DictItemItem, 0, len(d.Items))
		for _, di := range d.Items {
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
		result[code] = types.DictGetResp{Code: d.Code, Items: items}
	}

	return &types.DictBatchGetResp{Dicts: result}, nil
}
