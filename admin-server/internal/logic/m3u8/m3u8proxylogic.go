// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package m3u8

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	// m3u8HttpClient：用于拉取 m3u8 清单（小文件），短超时即可
	m3u8HttpClient = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	// mediaHttpClient：用于拉取 ts/分片等媒体流
	// 注意：这里不能用全局 Timeout（会在固定时间强制中断 io.Copy），
	// 由 Request Context 控制整体超时更可靠（支持边下边播/长连接传输）。
	mediaHttpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			// 仅限制“响应头”返回时间，避免上游卡死；实体传输用 ctx 控制
			ResponseHeaderTimeout: 15 * time.Second,
		},
	}
)

type M3u8ProxyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewM3u8ProxyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *M3u8ProxyLogic {
	return &M3u8ProxyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *M3u8ProxyLogic) M3u8Proxy(w http.ResponseWriter, r *http.Request, req *types.M3u8ProxyReq) error {
	targetURL := req.Url
	if targetURL == "" {
		return errs.New(errs.CodeBadRequest, "缺少url参数")
	}

	// 判断文件类型
	isM3U8 := strings.HasSuffix(strings.ToLower(targetURL), ".m3u8")

	// m3u8 文件不需要 Range 请求支持
	if isM3U8 {
		return l.handleM3U8Request(w, r, targetURL)
	}

	// 媒体文件支持 Range 请求（边下边播）
	return l.handleMediaFileWithRange(w, r, targetURL)
}

// handleM3U8Request 处理 m3u8 文件请求
func (l *M3u8ProxyLogic) handleM3U8Request(w http.ResponseWriter, r *http.Request, targetURL string) error {
	// 创建请求
	reqCtx, cancel := context.WithTimeout(l.ctx, 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodGet, targetURL, nil)
	if err != nil {
		return errs.Wrap(errs.CodeBadRequest, "创建请求失败", err)
	}

	// 发送请求
	resp, err := m3u8HttpClient.Do(httpReq)
	if err != nil {
		l.Errorf("请求目标失败: %v", err)
		return errs.New(errs.CodeBadGateway, "请求目标地址失败")
	}
	if resp == nil {
		return errs.New(errs.CodeBadGateway, "请求目标地址失败：响应为空")
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		l.Errorf("目标服务器返回错误状态: %d", resp.StatusCode)
		return errs.New(errs.CodeBadGateway, fmt.Sprintf("目标服务器返回错误: %d", resp.StatusCode))
	}

	return l.handleM3U8(w, resp.Body, targetURL, r)
}

// handleM3U8 处理 m3u8 文件，逐行处理避免全量读取
func (l *M3u8ProxyLogic) handleM3U8(w http.ResponseWriter, body io.Reader, targetURL string, r *http.Request) error {
	// 设置响应头
	// 注意：CORS 响应头已由 nginx 统一设置（/gateway/ location 中使用 always 参数）
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")

	// 获取基础 URL（用于补全相对路径）
	baseURL := targetURL[:strings.LastIndex(targetURL, "/")+1]
	domain := l.getDomain(r)

	// 使用 bufio.Scanner 逐行处理，避免全量读取
	scanner := bufio.NewScanner(body)
	var builder strings.Builder
	// 预分配容量，减少内存重新分配
	builder.Grow(4096)

	for scanner.Scan() {
		originalLine := scanner.Text()
		line := strings.TrimSpace(originalLine)

		// 保留空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			builder.WriteString(originalLine)
			builder.WriteString("\n")
			continue
		}

		// 处理媒体文件路径
		mediaURL := line
		if !strings.HasPrefix(mediaURL, "http") {
			mediaURL = baseURL + mediaURL
		}

		// 转换为代理地址
		proxyURL := fmt.Sprintf("%s/api/v1/m3u8/proxy?url=%s", domain, url.QueryEscape(mediaURL))
		builder.WriteString(proxyURL)
		builder.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		l.Errorf("读取 m3u8 失败: %v", err)
		return errs.Wrap(errs.CodeInternalError, "读取m3u8失败", err)
	}

	// 写入响应
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(builder.String())); err != nil {
		l.Errorf("写入响应失败: %v", err)
		return errs.Wrap(errs.CodeInternalError, "写入响应失败", err)
	}

	return nil
}

// handleMediaFileWithRange 处理 ts 等媒体文件，支持 Range 请求（边下边播）
func (l *M3u8ProxyLogic) handleMediaFileWithRange(w http.ResponseWriter, r *http.Request, targetURL string) error {
	// 创建请求（支持 Range 请求）
	reqCtx, cancel := context.WithTimeout(l.ctx, 120*time.Second) // 媒体文件允许更长时间（避免大分片/弱网 30s 超时）
	defer cancel()

	httpReq, err := http.NewRequestWithContext(reqCtx, http.MethodGet, targetURL, nil)
	if err != nil {
		return errs.Wrap(errs.CodeBadRequest, "创建请求失败", err)
	}

	// 透传客户端的 Range 请求头（关键：支持边下边播）
	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		httpReq.Header.Set("Range", rangeHeader)
	}

	// 发送请求
	resp, err := mediaHttpClient.Do(httpReq)
	if err != nil {
		l.Errorf("请求目标失败: %v", err)
		return errs.New(errs.CodeBadGateway, "请求目标地址失败")
	}
	if resp == nil {
		return errs.New(errs.CodeBadGateway, "请求目标地址失败：响应为空")
	}
	defer resp.Body.Close()

	// 检查响应状态（支持 200 和 206 Partial Content）
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		l.Errorf("目标服务器返回错误状态: %d", resp.StatusCode)
		return errs.New(errs.CodeBadGateway, fmt.Sprintf("目标服务器返回错误: %d", resp.StatusCode))
	}

	return l.handleMediaFile(w, resp)
}

// handleMediaFile 处理 ts 等媒体文件，直接透传（支持 Range 响应）
func (l *M3u8ProxyLogic) handleMediaFile(w http.ResponseWriter, resp *http.Response) error {
	if resp == nil {
		return errs.New(errs.CodeBadGateway, "响应为空")
	}
	// 需要排除的头部
	excludedHeaders := map[string]bool{
		"access-control-allow-origin":      true, // CORS 头由 nginx 统一设置
		"access-control-allow-methods":     true,
		"access-control-allow-headers":     true,
		"access-control-expose-headers":    true,
		"access-control-allow-credentials": true,
		"transfer-encoding":                true,
	}

	// 复制响应头（排除特定头部）
	for k, v := range resp.Header {
		lowerKey := strings.ToLower(k)
		if !excludedHeaders[lowerKey] {
			// 重要：保留 Range 请求相关的头部
			if lowerKey == "content-range" || lowerKey == "accept-ranges" || lowerKey == "content-length" {
				// 使用 Set 而不是 Add，确保只有一个值
				if len(v) > 0 {
					w.Header().Set(k, v[0])
				}
			} else {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}
		}
	}

	// 注意：CORS 响应头已由 nginx 统一设置（/gateway/ location 中使用 always 参数）
	// 添加 Accept-Ranges 头，告知客户端支持 Range 请求
	w.Header().Set("Accept-Ranges", "bytes")

	// 写入状态码（200 或 206 Partial Content）
	w.WriteHeader(resp.StatusCode)

	// 流式传输数据（边下边播的关键）
	if _, err := io.Copy(w, resp.Body); err != nil {
		l.Errorf("复制响应体失败: %v", err)
		return errs.Wrap(errs.CodeInternalError, "复制响应体失败", err)
	}

	return nil
}

// getDomain 获取请求的域名（处理代理头）
func (l *M3u8ProxyLogic) getDomain(r *http.Request) string {
	// 获取协议
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	// 获取 Host
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}

	// 获取前缀
	prefix := r.Header.Get("X-Forwarded-Prefix")

	return scheme + "://" + host + prefix
}
