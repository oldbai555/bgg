// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

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

	repo := sdkrepo.NewSdkAdminRepository(l.svcCtx.Repository)
	// status == 0 表示不按状态过滤，非0才过滤
	statusFilter := req.Status
	list, total, err := repo.ListInterfaces(l.ctx, req.Page, req.PageSize, req.Name, req.ApiCode, statusFilter)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询接口列表失败", err)
	}

	items := make([]types.SdkInterfaceItem, 0, len(list))
	for _, v := range list {
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
		Total: total,
		List:  items,
	}, nil
}
