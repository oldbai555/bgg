// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiListLogic {
	return &ApiListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiListLogic) ApiList(req *types.ApiListReq) (resp *types.ApiListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.ApiList(l.ctx, &iamclient.ApiListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询接口列表失败", err)
	}

	items := make([]types.ApiItem, 0, len(rpcResp.List))
	for _, a := range rpcResp.List {
		items = append(items, types.ApiItem{
			Id:          a.Id,
			Name:        a.Name,
			Method:      a.Method,
			Path:        a.Path,
			Description: a.Description,
			Status:      a.Status,
			CreatedAt:   a.CreatedAt,
		})
	}

	return &types.ApiListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
