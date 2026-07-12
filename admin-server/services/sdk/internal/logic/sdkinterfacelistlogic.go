package logic

import (
	"context"

	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkInterfaceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkInterfaceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkInterfaceListLogic {
	return &SdkInterfaceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkInterfaceListLogic) SdkInterfaceList(in *sdk.SdkInterfaceListRequest) (*sdk.SdkInterfaceListResponse, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.Admin.ListInterfaces(l.ctx, page, pageSize, in.Name, in.ApiCode, in.Status)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	items := make([]*sdk.SdkInterfaceItem, 0, len(list))
	for _, v := range list {
		items = append(items, &sdk.SdkInterfaceItem{
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

	return &sdk.SdkInterfaceListResponse{
		Total: total,
		List:  items,
	}, nil
}
