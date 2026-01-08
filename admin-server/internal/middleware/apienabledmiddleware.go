package middleware

import (
	"net/http"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/response"
)

// ApiEnabledMiddleware 接口启用校验中间件
// 只负责校验接口是否存在且已启用，不做登录和权限校验。
// 适用于 public_video 等无需登录但需要通过 AdminApi 配置开关的公共接口。
type ApiEnabledMiddleware struct {
	svcCtx *svc.ServiceContext
}

func NewApiEnabledMiddleware(svcCtx *svc.ServiceContext) *ApiEnabledMiddleware {
	return &ApiEnabledMiddleware{svcCtx: svcCtx}
}

func (m *ApiEnabledMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path

		apiRepo := repository.NewApiRepository(m.svcCtx.Repository)
		api, err := apiRepo.FindByMethodAndPath(r.Context(), method, path)
		if err != nil {
			// 接口不存在，直接返回 404
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeNotFound, "接口不存在"))
			return
		}

		// 接口未启用，返回禁止访问
		if api.Status != 1 {
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeForbidden, "接口未启用"))
			return
		}

		next(w, r)
	}
}
