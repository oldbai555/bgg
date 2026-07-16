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

type BlogArticleAuditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleAuditLogic {
	return &BlogArticleAuditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BlogArticleAudit 薄胶水：解析 HTTP 请求（含从 JWT context 取当前登录审核员）-> 调
// ContentRPC -> 映射响应，实际业务逻辑（含审计日志回调 IamCallback.RecordAuditLog）已经
// 搬进 services/content/internal/logic/blogarticleauditlogic.go。
func (l *BlogArticleAuditLogic) BlogArticleAudit(req *types.BlogArticleAuditReq) (resp *types.Response, err error) {
	var operatorUserID uint64
	var operatorUsername string
	if u, ok := jwthelper.FromContext(l.ctx); ok {
		operatorUserID = u.UserID
		operatorUsername = u.Username
	}

	_, err = l.svcCtx.ContentRPC.BlogArticleAudit(l.ctx, &contentclient.BlogArticleAuditRequest{
		Id:               req.Id,
		Result:           req.Result,
		Remark:           req.Remark,
		OperatorUserId:   operatorUserID,
		OperatorUsername: operatorUsername,
	})
	if err != nil {
		return nil, errs.WrapGRPCError("审核失败", err)
	}

	return &types.Response{Code: int(errs.CodeOK), Message: "审核成功"}, nil
}
