// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package sdk

import (
	"context"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/sdkclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkInterfaceListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkInterfaceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceListLogic {
	return &SdkInterfaceListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkInterfaceListLogic) SdkInterfaceList(req *types.SdkInterfaceListReq) (resp *types.SdkInterfaceListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	rpcResp, err := l.svcCtx.SdkRPC.SdkInterfaceList(l.ctx, &sdkclient.SdkInterfaceListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		ApiCode:  req.ApiCode,
		Status:   req.Status,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询接口列表失败", err)
	}

	items := make([]types.SdkInterfaceItem, 0, len(rpcResp.List))
	for _, v := range rpcResp.List {
		items = append(items, types.SdkInterfaceItem{
			Id:               v.Id,
			Name:             v.Name,
			ApiCode:          v.ApiCode,
			Path:             v.Path,
			Method:           v.Method,
			RateLimitDefault: v.RateLimitDefault,
			Status:           v.Status,
			Remark:           v.Remark,
			CreatedAt:        v.CreatedAt,
		})
	}

	return &types.SdkInterfaceListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
