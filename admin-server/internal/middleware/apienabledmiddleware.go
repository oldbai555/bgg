package middleware

import (
	"net/http"
	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/repository"
	iamrepo "postapocgame/admin-server/internal/repository/iam"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/response"
	"strings"
)

// ApiEnabledMiddleware 接口启用校验中间件
// 只负责校验接口是否存在且已启用，不做登录和权限校验。
// 适用于 public_video 等无需登录但需要通过 AdminApi 配置开关的公共接口。
type ApiEnabledMiddleware struct {
	repo *repository.Repository
}

func NewApiEnabledMiddleware(repo *repository.Repository) *ApiEnabledMiddleware {
	return &ApiEnabledMiddleware{repo: repo}
}

func (m *ApiEnabledMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path

		apiRepo := iamrepo.NewApiRepository(m.repo)
		api, err := apiRepo.FindByMethodAndPath(r.Context(), method, path)
		if err != nil {
			// 接口不存在，直接返回 404
			// 对于需要 CORS 的接口（如 m3u8、video_collect），设置 CORS 头
			m.setCORSIfNeeded(w, path)
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeNotFound, "接口不存在"))
			return
		}

		// 接口未启用，返回禁止访问
		if api.Status != consts.Open {
			// 对于需要 CORS 的接口（如 m3u8、video_collect），设置 CORS 头
			m.setCORSIfNeeded(w, path)
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeForbidden, "接口未启用"))
			return
		}

		next(w, r)
	}
}

// setCORSIfNeeded 对于需要 CORS 的接口，设置 CORS 响应头
func (m *ApiEnabledMiddleware) setCORSIfNeeded(w http.ResponseWriter, path string) {
	// 判断是否是需要 CORS 的接口
	if strings.Contains(path, "/m3u8/proxy") || strings.Contains(path, "/videos/collect") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Range,X-Requested-With")
		if strings.Contains(path, "/m3u8/proxy") {
			w.Header().Set("Access-Control-Expose-Headers", "Content-Range,Accept-Ranges,Content-Length")
		}
		w.Header().Set("Access-Control-Max-Age", "86400") // 24小时
	}
}
