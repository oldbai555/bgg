package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/sdk/internal/svc"
	"postapocgame/admin-server/services/sdk/sdk"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkCallLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkCallLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkCallLogListLogic {
	return &SdkCallLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SdkCallLogListLogic) SdkCallLogList(in *sdk.SdkCallLogListRequest) (*sdk.SdkCallLogListResponse, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, err := l.svcCtx.Admin.ListCallLogs(l.ctx, page, pageSize, in.SdkKeyId, in.ApiCode, in.RespCode, in.Ip, in.StartTime, in.EndTime)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询调用记录失败", err))
	}

	items := make([]*sdk.SdkCallLogItem, 0, len(list))
	for _, v := range list {
		items = append(items, &sdk.SdkCallLogItem{
			Id:             v.Id,
			SdkKeyId:       v.SdkKeyId,
			SdkInterfaceId: v.SdkInterfaceId,
			ApiCode:        v.ApiCode,
			Path:           v.Path,
			Method:         v.Method,
			Ip:             v.Ip,
			RespCode:       v.RespCode,
			DurationMs:     v.DurationMs,
			CreatedAt:      v.CreatedAt,
		})
	}

	return &sdk.SdkCallLogListResponse{
		Total: total,
		List:  items,
	}, nil
}
