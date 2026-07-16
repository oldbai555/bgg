package logic

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogSocialInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogSocialInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoDeleteLogic {
	return &BlogSocialInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogSocialInfoDelete 迁移自 internal/logic/blog/social_info/blog_social_info_delete_logic.go。
func (l *BlogSocialInfoDeleteLogic) BlogSocialInfoDelete(in *content.BlogSocialInfoDeleteRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "社交信息ID不能为空"))
	}
	info, err := l.svcCtx.SocialInfo.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询社交信息失败", err))
	}
	if info == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "社交信息不存在"))
	}
	if err := l.svcCtx.SocialInfo.Delete(l.ctx, in.Id); err != nil {
		return nil, toGRPCStatus(err)
	}
	return &content.Empty{}, nil
}
