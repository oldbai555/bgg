// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package social_info

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

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

func (l *BlogSocialInfoDeleteLogic) BlogSocialInfoDelete(req *types.BlogSocialInfoDeleteReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "社交信息ID不能为空")
	}

	// 检查记录是否存在
	info, err := l.svcCtx.Domain.Blog.SocialInfo.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询社交信息失败", err)
	}
	if info == nil {
		return nil, errs.New(errs.CodeNotFound, "社交信息不存在")
	}

	// 执行软删除
	if err := l.svcCtx.Domain.Blog.SocialInfo.Delete(l.ctx, req.Id); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "删除成功",
	}, nil
}
