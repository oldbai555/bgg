package middleware

import (
	"net/http"

	iamdomain "postapocgame/admin-server/internal/domain/iam"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"
	"postapocgame/admin-server/pkg/response"
)

// PermissionMiddleware 权限鉴权中间件
type PermissionMiddleware struct {
	resolver *iamdomain.PermissionResolver
}

func NewPermissionMiddleware(resolver *iamdomain.PermissionResolver) *PermissionMiddleware {
	return &PermissionMiddleware{resolver: resolver}
}

func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := jwthelper.FromContext(r.Context())
		if !ok {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
			return
		}

		allowed, err := m.resolver.CanAccess(r.Context(), user.UserID, r.Method, r.URL.Path)
		if err != nil {
			response.ErrorCtx(r.Context(), w, err)
			return
		}
		if !allowed {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeForbidden, "无权限访问该接口"))
			return
		}

		next(w, r)
	}
}
