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

type PublicBlogAuthorInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublicBlogAuthorInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublicBlogAuthorInfoLogic {
	return &PublicBlogAuthorInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PublicBlogAuthorInfo 薄胶水：原来的跨域读取 IAM 用户信息（TODO(phase2-content-rpc)）
// 已经在 content-rpc 侧改成回调 IamCallback.GetUserProfile，实际业务逻辑搬进
// services/content/internal/logic/publicblogauthorinfologic.go。
func (l *PublicBlogAuthorInfoLogic) PublicBlogAuthorInfo() (resp *types.PublicBlogAuthorInfoResp, err error) {
	rpcResp, err := l.svcCtx.ContentRPC.PublicBlogAuthorInfo(l.ctx, &contentclient.PublicBlogGlobalRequest{})
	if err != nil {
		return nil, errs.WrapGRPCError("查询作者信息失败", err)
	}
	return &types.PublicBlogAuthorInfoResp{
		Id:        rpcResp.Id,
		Nickname:  rpcResp.Nickname,
		Avatar:    rpcResp.Avatar,
		Signature: rpcResp.Signature,
	}, nil
}
