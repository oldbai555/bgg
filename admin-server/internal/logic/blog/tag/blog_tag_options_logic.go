package tag

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

// BlogTagOptionsLogic 标签下拉选项逻辑
type BlogTagOptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagOptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagOptionsLogic {
	return &BlogTagOptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogTagOptions 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogtagoptionslogic.go。
func (l *BlogTagOptionsLogic) BlogTagOptions(req *types.BlogTagOptionsReq) (resp *types.BlogTagOptionsResp, err error) {
	limit := int64(0)
	if req != nil {
		limit = req.Limit
	}
	rpcResp, err := l.svcCtx.ContentRPC.BlogTagOptions(l.ctx, &contentclient.BlogTagOptionsRequest{Limit: limit})
	if err != nil {
		return nil, errs.WrapGRPCError("查询标签选项失败", err)
	}
	items := make([]types.BlogTagOptionItem, 0, len(rpcResp.List))
	for _, t := range rpcResp.List {
		items = append(items, types.BlogTagOptionItem{Id: t.Id, Name: t.Name})
	}
	return &types.BlogTagOptionsResp{List: items}, nil
}
