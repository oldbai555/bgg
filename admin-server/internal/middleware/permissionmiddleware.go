package middleware

import (
	"net/http"

	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/pkg/response"
	"postapocgame/admin-server/services/iam/iamclient"

	"github.com/zeromicro/go-zero/core/logx"
)

// PermissionMiddleware 权限鉴权中间件。iam 域拆分成独立服务后，原来直接持有的
// *iamdomain.PermissionResolver 换成 iam-rpc 的 zrpc client——这是这个中间件第一次
// 真正切换成调 RPC，见 18-service-extraction-runbook.md 2.5 节。
type PermissionMiddleware struct {
	iamRPC iamclient.Iam
}

func NewPermissionMiddleware(iamRPC iamclient.Iam) *PermissionMiddleware {
	return &PermissionMiddleware{iamRPC: iamRPC}
}

func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := jwthelper.FromContext(r.Context())
		if !ok {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
			return
		}

		resp, err := m.iamRPC.CheckPermission(r.Context(), &iamclient.CheckPermissionRequest{
			UserId: user.UserID,
			Method: r.Method,
			Path:   r.URL.Path,
		})
		if err != nil {
			response.ErrorCtx(r.Context(), w, errs.WrapGRPCError("权限校验失败", err))
			return
		}
		if !resp.Allowed {
			logx.Infof("权限校验拒绝: userId=%d, method=%s, path=%s, reason=%s", user.UserID, r.Method, r.URL.Path, resp.Reason)
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeForbidden, "无权限访问该接口"))
			return
		}

		next(w, r)
	}
}
