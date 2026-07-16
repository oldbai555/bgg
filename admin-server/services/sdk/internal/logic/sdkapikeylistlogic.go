package logic

import (
	"context"

	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkApiKeyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkApiKeyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkApiKeyListLogic {
	return &SdkApiKeyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkApiKeyListLogic) SdkApiKeyList(in *sdk.SdkApiKeyListRequest) (*sdk.SdkApiKeyListResponse, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.Admin.ListSdkKeys(l.ctx, page, pageSize, in.Name, in.Status)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	items := make([]*sdk.SdkApiKeyItem, 0, len(list))
	for _, k := range list {
		items = append(items, &sdk.SdkApiKeyItem{
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

	return &sdk.SdkApiKeyListResponse{
		Total: total,
		List:  items,
	}, nil
}
