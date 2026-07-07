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

func (l *SdkApiKeyListLogic) SdkApiKeyList(req *types.SdkApiKeyListReq) (resp *types.SdkApiKeyListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, 20, 100)

	repo := sdkrepo.NewSdkAdminRepository(l.svcCtx.Repository)
	// status == 0 表示不按状态过滤，非0才过滤
	statusFilter := req.Status
	list, total, err := repo.ListSdkKeys(l.ctx, req.Page, req.PageSize, req.Name, statusFilter)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询 API Key 列表失败", err)
	}

	items := make([]types.SdkApiKeyItem, 0, len(list))
	for _, k := range list {
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
		Total: total,
		List:  items,
	}, nil
}
