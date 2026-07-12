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

type SdkApiKeyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkApiKeyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyListLogic {
	return &SdkApiKeyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SdkApiKeyList 薄胶水：解析 HTTP 请求 -> 拼一次 SdkRPC 请求 -> 映射响应，sdk 域的实际
// 业务逻辑已经搬进 services/sdk/internal/logic/sdkapikeylistlogic.go。
func (l *SdkApiKeyListLogic) SdkApiKeyList(req *types.SdkApiKeyListReq) (resp *types.SdkApiKeyListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	rpcResp, err := l.svcCtx.SdkRPC.SdkApiKeyList(l.ctx, &sdkclient.SdkApiKeyListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		Status:   req.Status,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询 API Key 列表失败", err)
	}

	items := make([]types.SdkApiKeyItem, 0, len(rpcResp.List))
	for _, k := range rpcResp.List {
		items = append(items, types.SdkApiKeyItem{
			Id:          k.Id,
			Name:        k.Name,
			ApiKey:      k.ApiKey,
			ApiSecret:   k.ApiSecret,
			Status:      k.Status,
			ExpireAt:    k.ExpireAt,
			IpWhitelist: k.IpWhitelist,
			Remark:      k.Remark,
			CreatedAt:   k.CreatedAt,
		})
	}

	return &types.SdkApiKeyListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
