// 自定义路由注册文件
// 此文件不会被 goctl 自动生成覆盖，用于注册 WebSocket 等自定义路由

package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"postapocgame/admin-server/internal/consts"
	chat "postapocgame/admin-server/internal/handler/chat"
	"postapocgame/admin-server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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

	// 静态文件服务：/api/v1/files/uploads/* -> ./uploads/*
	// 用于访问上传的文件，不需要认证（公开访问）
	// go-zero 不支持通配符路由，使用自定义 Handler 处理所有匹配的请求
	fileServerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查是否是文件上传路径
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		// 从路径中提取文件名（去掉前缀 /api/v1/files/uploads/）
		path := r.URL.Path
		if !strings.HasPrefix(path, consts.PathFileUploads+"/") {
			http.NotFound(w, r)
			return
		}

		path = strings.TrimPrefix(path, consts.PathFileUploads+"/")

		if path == "" || strings.Contains(path, "..") {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		// 构建文件路径
		filePath := filepath.Join(consts.UploadDir, path)

		// 检查文件是否存在
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			logx.Errorf("获取文件绝对路径失败: %v", err)
			http.NotFound(w, r)
			return
		}

		// 确保文件路径在 uploads 目录内（防止路径遍历）
		uploadAbsPath, _ := filepath.Abs(consts.UploadDir)
		if !strings.HasPrefix(absPath, uploadAbsPath) {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// 文件存在且路径安全，直接服务文件
		http.ServeFile(w, r, filePath)
	})

	// 使用 server.AddRoute 直接注册完整路径
	// 注意：需要在 RegisterHandlers 之后注册，以确保优先级
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    consts.PathFileUploads + "/:filename",
		Handler: fileServerHandler,
	})

	logx.Infof("静态文件服务已注册: %s/* -> %s", consts.PathFileUploads, consts.UploadDir)

	// 注意：博客标签下拉选项接口 /blog/tags/options 已由 goctl 自动生成到 routes.go
	// 无需在此重复注册
}
