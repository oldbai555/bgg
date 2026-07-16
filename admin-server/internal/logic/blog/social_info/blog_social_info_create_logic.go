// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package social_info

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogSocialInfoCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogSocialInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoCreateLogic {
	return &BlogSocialInfoCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogSocialInfoCreate 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogsocialinfocreatelogic.go。
func (l *BlogSocialInfoCreateLogic) BlogSocialInfoCreate(req *types.BlogSocialInfoCreateReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogSocialInfoCreate(l.ctx, &contentclient.BlogSocialInfoCreateRequest{
		Name:     req.Name,
		Url:      req.Url,
		Remark:   req.Remark,
		Status:   req.Status,
		OrderNum: req.OrderNum,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建社交信息失败", err)
	}
	return &types.Response{Code: 0, Message: "创建成功"}, nil
}
