package middleware

import (
	"net/http"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/pkg/response"
	"postapocgame/admin-server/services/iam/iamclient"
)

// ApiEnabledMiddleware 接口启用校验中间件
// 只负责校验接口是否存在且已启用，不做登录和权限校验。
// 适用于 public_video 等无需登录但需要通过 AdminApi 配置开关的公共接口。
type ApiEnabledMiddleware struct {
	iamRPC iamclient.Iam
}

func NewApiEnabledMiddleware(iamRPC iamclient.Iam) *ApiEnabledMiddleware {
	return &ApiEnabledMiddleware{iamRPC: iamRPC}
}

func (m *ApiEnabledMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path

		resp, err := m.iamRPC.CheckApiEnabled(r.Context(), &iamclient.CheckApiEnabledRequest{Method: method, Path: path})
		if err != nil {
			m.setCORSIfNeeded(w, path)
			response.ErrorCtx(r.Context(), w, errs.WrapGRPCError("检查接口状态失败", err))
			return
		}
		if !resp.Exists {
			m.setCORSIfNeeded(w, path)
			response.ErrorCtx(r.Context(), w, errs.New(errs.CodeNotFound, "接口不存在"))
			return
		}
		if !resp.Enabled {
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
	if strings.Contains(path, consts.CORSPathM3U8Proxy) || strings.Contains(path, consts.CORSPathVideoCollect) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,Range,X-Requested-With")
		if strings.Contains(path, consts.CORSPathM3U8Proxy) {
			w.Header().Set("Access-Control-Expose-Headers", "Content-Range,Accept-Ranges,Content-Length")
		}
		w.Header().Set("Access-Control-Max-Age", "86400") // 24小时
	}
}
