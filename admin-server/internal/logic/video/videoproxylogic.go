// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package video

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoProxyLogic {
	return &VideoProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoProxyLogic) VideoProxy(w http.ResponseWriter, r *http.Request, req *types.VideoProxyReq) error {
	if req == nil || req.Url == "" {
		return errs.New(errs.CodeBadRequest, "视频URL不能为空")
	}

	// 从字典中获取视频代理地址配置
	dictTypeRepo := repository.NewDictTypeRepository(l.svcCtx.Repository)
	dictType, err := dictTypeRepo.FindByCode(l.ctx, "video_proxy_url")
	if err != nil {
		l.Errorf("获取视频代理配置失败: %v", err)
		// 如果字典配置不存在，使用默认代理地址
		return l.proxyVideo(w, r, req.Url, "")
	}

	dictItemRepo := repository.NewDictItemRepository(l.svcCtx.Repository)
	items, err := dictItemRepo.FindByTypeID(l.ctx, dictType.Id)
	if err != nil || len(items) == 0 {
		l.Errorf("获取视频代理配置项失败: %v", err)
		// 如果配置项不存在，使用默认代理地址
		return l.proxyVideo(w, r, req.Url, "")
	}

	// 使用第一个有效的代理地址
	proxyBaseUrl := items[0].Value
	return l.proxyVideo(w, r, req.Url, proxyBaseUrl)
}

func (l *VideoProxyLogic) proxyVideo(w http.ResponseWriter, r *http.Request, videoUrl string, proxyBaseUrl string) error {
	// 如果提供了代理基础URL，构建代理URL
	var targetUrl string
	if proxyBaseUrl != "" {
		// 将视频URL编码后作为参数传递
		encodedUrl := url.QueryEscape(videoUrl)
		targetUrl = fmt.Sprintf("%s?url=%s", proxyBaseUrl, encodedUrl)
	} else {
		// 直接使用视频URL
		targetUrl = videoUrl
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建代理请求
	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetUrl, r.Body)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, "创建代理请求失败", err)
	}

	// 复制请求头（排除一些不需要的头部）
	for key, values := range r.Header {
		// 排除Host、Connection等头部
		if strings.ToLower(key) == "host" || strings.ToLower(key) == "connection" {
			continue
		}
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// 设置Referer为原始视频URL（某些视频服务器需要）
	proxyReq.Header.Set("Referer", videoUrl)

	// 执行代理请求
	resp, err := client.Do(proxyReq)
	if err != nil {
		return errs.Wrap(errs.CodeInternalError, "代理请求失败", err)
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		// 排除一些不需要的头部
		if strings.ToLower(key) == "content-length" {
			// Content-Length由Go自动设置
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		l.Errorf("复制响应体失败: %v", err)
		return errs.Wrap(errs.CodeInternalError, "复制响应体失败", err)
	}

	return nil
}
