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

type SdkCallLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSdkCallLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkCallLogListLogic {
	return &SdkCallLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SdkCallLogListLogic) SdkCallLogList(req *types.SdkCallLogListReq) (resp *types.SdkCallLogListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 200)

	rpcResp, err := l.svcCtx.SdkRPC.SdkCallLogList(l.ctx, &sdkclient.SdkCallLogListRequest{
		Page:      req.Page,
		PageSize:  req.PageSize,
		SdkKeyId:  req.SdkKeyId,
		ApiCode:   req.ApiCode,
		RespCode:  req.RespCode,
		Ip:        req.Ip,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询调用记录失败", err)
	}

	items := make([]types.SdkCallLogItem, 0, len(rpcResp.List))
	for _, v := range rpcResp.List {
		items = append(items, types.SdkCallLogItem{
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

	return &types.SdkCallLogListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
