// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article_audit

import (
	"context"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/services/content/contentclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleAuditUnpublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleAuditUnpublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleAuditUnpublishLogic {
	return &BlogArticleAuditUnpublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleAuditUnpublish 薄胶水，实际业务逻辑已经搬进
// services/content/internal/logic/blogarticleauditunpublishlogic.go。
func (l *BlogArticleAuditUnpublishLogic) BlogArticleAuditUnpublish(req *types.BlogArticleAuditUnpublishReq) (resp *types.Response, err error) {
	var operatorUserID uint64
	var operatorUsername string
	if u, ok := jwthelper.FromContext(l.ctx); ok {
		operatorUserID = u.UserID
		operatorUsername = u.Username
	}

	_, err = l.svcCtx.ContentRPC.BlogArticleAuditUnpublish(l.ctx, &contentclient.BlogArticleAuditUnpublishRequest{
		Id:               req.Id,
		Remark:           req.Remark,
		OperatorUserId:   operatorUserID,
		OperatorUsername: operatorUsername,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("下架失败", err)
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "下架成功"}, nil
}
