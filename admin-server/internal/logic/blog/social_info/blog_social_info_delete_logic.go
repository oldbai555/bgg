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

type BlogSocialInfoDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogSocialInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoDeleteLogic {
	return &BlogSocialInfoDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogSocialInfoDelete 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogsocialinfodeletelogic.go。
func (l *BlogSocialInfoDeleteLogic) BlogSocialInfoDelete(req *types.BlogSocialInfoDeleteReq) (resp *types.Response, err error) {
	_, err = l.svcCtx.ContentRPC.BlogSocialInfoDelete(l.ctx, &contentclient.BlogSocialInfoDeleteRequest{Id: req.Id})
	if err != nil {
		return nil, errs.WrapGRPCError("删除社交信息失败", err)
	}
	return &types.Response{Code: 0, Message: "删除成功"}, nil
}
