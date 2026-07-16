// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigListLogic {
	return &ConfigListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigListLogic) ConfigList(req *types.ConfigListReq) (resp *types.ConfigListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	rpcResp, err := l.svcCtx.IamRPC.ConfigList(l.ctx, &iamclient.ConfigListRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
		Group:    req.Group,
		Key:      req.Key,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询配置列表失败", err)
	}

	items := make([]types.ConfigItem, 0, len(rpcResp.List))
	for _, c := range rpcResp.List {
		items = append(items, types.ConfigItem{
			Id:          c.Id,
			Group:       c.Group,
			Key:         c.Key,
			Value:       c.Value,
			ConfigType:  c.ConfigType,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
		})
	}

	return &types.ConfigListResp{
		Total: rpcResp.Total,
		List:  items,
	}, nil
}
