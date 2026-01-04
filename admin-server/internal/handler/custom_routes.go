// 自定义路由注册文件
// 此文件不会被 goctl 自动生成覆盖，用于注册 WebSocket 等自定义路由

package handler

import (
	"net/http"

	"postapocgame/admin-server/internal/consts"
	chat "postapocgame/admin-server/internal/handler/chat"
	"postapocgame/admin-server/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// RegisterCustomRoutes 注册自定义路由（WebSocket、静态文件服务等）
// 此函数应在 RegisterHandlers 之后调用
func RegisterCustomRoutes(server *rest.Server, serverCtx *svc.ServiceContext) {
	// WebSocket 路由（不需要权限中间件，在 Handler 内部验证）
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    consts.PathChatWS,
		Handler: chat.ChatWSHandler(serverCtx),
	})

	// 静态文件服务路由（用于直接访问上传的文件，通过nginx代理）
	// 路径：/api/v1/uploads/* -> ./uploads/*
	// nginx配置：/files/uploads/* -> /api/v1/uploads/*
	server.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/uploads/*",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/api/v1/uploads/", http.FileServer(http.Dir("./uploads"))).ServeHTTP(w, r)
		}),
	})

	// 注意：操作日志中间件需要在 routes.go 中手动添加到所有需要认证的路由组
	// 由于 routes.go 是自动生成的，每次执行 generate-api.sh 后需要手动添加
	// 在所有需要认证的路由组的 WithMiddlewares 中添加 serverCtx.OperationLogMiddleware
}
