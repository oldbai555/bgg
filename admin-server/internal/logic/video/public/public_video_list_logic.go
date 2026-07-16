// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicVideoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicVideoListLogic {
	return &PublicVideoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicVideoList 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/publicvideolistlogic.go。
func (l *PublicVideoListLogic) PublicVideoList(req *types.PublicVideoListReq) (resp *types.PublicVideoListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicVideoList(l.ctx, &contentclient.PublicVideoListRequest{
		Page:    req.Page,
		Size:    req.Size,
		Content: req.Content,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询视频列表失败", err)
	}

	items := make([]types.PublicVideoItem, 0, len(rpcResp.List))
	for _, v := range rpcResp.List {
		items = append(items, types.PublicVideoItem{Id: v.Id, Uuid: v.Uuid, Name: v.Name, GodNum: v.GodNum})
	}

	return &types.PublicVideoListResp{List: items, Page: rpcResp.Page, Size: rpcResp.Size, Total: rpcResp.Total}, nil
}
