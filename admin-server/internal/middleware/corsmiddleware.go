package middleware

import (
	"net/http"

	"postapocgame/admin-server/internal/consts"
)

// CorsMiddleware CORS 跨域中间件
// 用于处理所有 /public 接口的跨域请求
type CorsMiddleware struct{}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头
		origin := r.Header.Get("Origin")

		// 允许所有来源（生产环境可以根据需要配置白名单）
		// 如果 origin 存在，使用 origin；否则使用 *
		if origin != "" {
			w.Header().Set(consts.HeaderAccessControlAllowOrigin, origin)
			// 当设置了具体的 origin 时，可以允许携带凭证
			w.Header().Set(consts.HeaderAccessControlAllowCredentials, "true")
		} else {
			// 没有 origin 时，使用 * 但不允许携带凭证（浏览器限制）
			w.Header().Set(consts.HeaderAccessControlAllowOrigin, consts.CORSAllowOriginAll)
		}

		// 允许的 HTTP 方法
		w.Header().Set(consts.HeaderAccessControlAllowMethods, consts.CORSAllowMethodsAll)

		// 允许的请求头
		w.Header().Set(consts.HeaderAccessControlAllowHeaders, consts.CORSAllowHeadersPublic)

		// 允许暴露的响应头
		w.Header().Set(consts.HeaderAccessControlExposeHeaders, consts.CORSExposeHeadersDefault)

		// 预检请求的缓存时间（24小时）
		w.Header().Set(consts.HeaderAccessControlMaxAge, consts.CORSMaxAge)

		// 处理 OPTIONS 预检请求
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 继续处理实际请求
		next(w, r)
	}
}
