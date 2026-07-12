package logic

import (
	"context"

	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublicBlogAuthorInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublicBlogAuthorInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogAuthorInfoLogic {
	return &PublicBlogAuthorInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PublicBlogAuthorInfo 迁移自 internal/logic/blog/public/public_blog_author_info_logic.go。
// 原实现直连 l.svcCtx.Domain.IAM.User.FindByID(l.ctx, 1)（跨域读 IAM 用户，标注
// TODO(phase2-content-rpc)），content-rpc 拆分后改成回调 IamCallback.GetUserProfile，
// 这是 18-service-extraction-runbook.md 2.4 节点名的唯一一处 content-rpc 现场需要判断的
// 跨服务点。
func (l *PublicBlogAuthorInfoLogic) PublicBlogAuthorInfo(in *content.PublicBlogGlobalRequest) (*content.PublicBlogAuthorInfoResponse, error) {
	profile, err := l.svcCtx.IamCallback.GetUserProfile(l.ctx, &iamcallbackpb.GetUserProfileRequest{UserId: 1})
	if err != nil || !profile.Exists {
		// 用户不存在/已删除/回调失败均返回默认值，和原实现的降级语义一致。
		return &content.PublicBlogAuthorInfoResponse{
			Id:       1,
			Nickname: "管理员",
		}, nil
	}

	return &content.PublicBlogAuthorInfoResponse{
		Id:        1,
		Nickname:  profile.Nickname,
		Avatar:    profile.Avatar,
		Signature: profile.Signature,
	}, nil
}
