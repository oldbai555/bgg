// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict_type

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type DictTypeListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDictTypeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DictTypeListLogic {
	return &DictTypeListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DictTypeListLogic) DictTypeList(req *types.DictTypeListReq) (resp *types.DictTypeListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.DictTypeList(l.ctx, &iamclient.DictTypeListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Code:     req.Code,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询字典类型列表失败", err)
	}

	items := make([]types.DictTypeItem, 0, len(rpcResp.List))
	for _, dt := range rpcResp.List {
		items = append(items, types.DictTypeItem{
			Id:          dt.Id,
			Name:        dt.Name,
			Code:        dt.Code,
			Description: dt.Description,
			Status:      dt.Status,
			CreatedAt:   dt.CreatedAt,
		})
	}

	return &types.DictTypeListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
