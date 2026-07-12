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

type BlogSocialInfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogSocialInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoUpdateLogic {
	return &BlogSocialInfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogSocialInfoUpdate 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogsocialinfoupdatelogic.go。
func (l *BlogSocialInfoUpdateLogic) BlogSocialInfoUpdate(req *types.BlogSocialInfoUpdateReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogSocialInfoUpdate(l.ctx, &contentclient.BlogSocialInfoUpdateRequest{
		Id:       req.Id,
		Name:     req.Name,
		Url:      req.Url,
		Remark:   req.Remark,
		Status:   req.Status,
		OrderNum: req.OrderNum,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("更新社交信息失败", err)
	}
	return &types.Response{Code: 0, Message: "更新成功"}, nil
}
