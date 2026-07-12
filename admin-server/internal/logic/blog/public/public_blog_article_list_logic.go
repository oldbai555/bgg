// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package public

import (
	"context"
	"strings"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogArticleListLogic {
	return &PublicBlogArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogArticleList 薄胶水：摘要截断长度的静态配置、Markdown 去除等实际业务逻辑已经
// 搬进 services/content/internal/logic/publicblogarticlelistlogic.go。
func (l *PublicBlogArticleListLogic) PublicBlogArticleList(req *types.PublicBlogArticleListReq) (resp *types.PublicBlogArticleListResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogArticleList(l.ctx, &contentclient.PublicBlogArticleListRequest{
		Page:    req.Page,
		Size:    req.Size,
		TagId:   req.TagId,
		Keyword: strings.TrimSpace(req.Keyword),
	})
	if err != nil {
		return nil, errs.WrapGRPCError("查询文章列表失败", err)
	}

	items := make([]types.PublicBlogArticleItem, 0, len(rpcResp.List))
	for _, a := range rpcResp.List {
		items = append(items, types.PublicBlogArticleItem{
			Id:          a.Id,
			Title:       a.Title,
			Cover:       a.Cover,
			AuthorName:  a.AuthorName,
			Summary:     a.Summary,
			TagNames:    a.TagNames,
			PublishTime: a.PublishTime,
			IsTop:       a.IsTop,
		})
	}

	return &types.PublicBlogArticleListResp{List: items, Page: rpcResp.Page, Size: rpcResp.Size, Total: rpcResp.Total}, nil
}
