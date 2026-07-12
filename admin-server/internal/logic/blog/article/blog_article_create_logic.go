// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleCreateLogic {
	return &BlogArticleCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleCreate 薄胶水：解析 HTTP 请求（含从 JWT context 取当前登录用户作为作者）
// -> 调 ContentRPC -> 映射响应，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticlecreatelogic.go。
func (l *BlogArticleCreateLogic) BlogArticleCreate(req *types.BlogArticleCreateReq) (resp *types.Response, err error) {
	var operatorUserID uint64
	var operatorUsername string
	if u, ok := jwthelper.FromContext(l.ctx); ok {
		operatorUserID = u.UserID
		operatorUsername = u.Username
	}

	_, err = l.svcCtx.ContentRPC.BlogArticleCreate(l.ctx, &contentclient.BlogArticleCreateRequest{
		Title:            req.Title,
		Content:          req.Content,
		TagIds:           req.TagIds,
		Cover:            req.Cover,
		Summary:          req.Summary,
		OperatorUserId:   operatorUserID,
		OperatorUsername: operatorUsername,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("创建文章失败", err)
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "创建成功"}, nil
}
