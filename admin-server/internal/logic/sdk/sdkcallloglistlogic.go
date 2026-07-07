// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package sdk

import (
	"context"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
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

	repo := sdkrepo.NewSdkAdminRepository(l.svcCtx.Repository)
	list, total, err := repo.ListCallLogs(l.ctx, req.Page, req.PageSize, req.SdkKeyId, req.ApiCode, req.RespCode, req.Ip, req.StartTime, req.EndTime)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询调用记录失败", err)
	}

	items := make([]types.SdkCallLogItem, 0, len(list))
	for _, v := range list {
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
		Total: total,
		List:  items,
	}, nil
}
