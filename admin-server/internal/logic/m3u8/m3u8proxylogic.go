// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package m3u8

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
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
	// 获取目标URL
	targetURL := req.Url
	if targetURL == "" {
		return errs.New(errs.CodeBadRequest, "缺少url参数")
	}

	// 获取请求的域名（处理代理头）
	domain := l.getDomain(r)

	// 向目标地址发送请求
	resp, err := http.Get(targetURL)
	if err != nil {
		l.Errorf("请求目标失败: %v", err)
		return errs.New(errs.CodeBadGateway, "请求目标地址失败")
	}
	defer resp.Body.Close()

	// 如果是 m3u8 文件
	if strings.HasSuffix(strings.ToLower(targetURL), ".m3u8") {
		// 先删除可能存在的 CORS 头（确保干净）
		w.Header().Del("Access-Control-Allow-Origin")
		w.Header().Del("Access-Control-Allow-Methods")
		w.Header().Del("Access-Control-Allow-Headers")
		w.Header().Del("Access-Control-Expose-Headers")
		w.Header().Del("Access-Control-Allow-Credentials")

		// 设置响应头（在写入响应之前）
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		// 设置CORS响应头（确保只设置一次）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			l.Errorf("读取 m3u8 失败: %v", err)
			return errs.New(errs.CodeInternalError, "读取m3u8失败")
		}

		m3u8Content := string(body)
		baseURL := targetURL[:strings.LastIndex(targetURL, "/")+1]

		lines := strings.Split(m3u8Content, "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				// 相对路径补全
				if !strings.HasPrefix(line, "http") {
					line = baseURL + line
				}
				// 转换为代理地址
				proxyURL := fmt.Sprintf("%s/api/v1/m3u8/proxy?url=%s", domain, url.QueryEscape(line))
				lines[i] = proxyURL
			}
		}

		newM3U8 := strings.Join(lines, "\n")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(newM3U8))
		if err != nil {
			l.Errorf("写入响应失败: %v", err)
			return errs.New(errs.CodeInternalError, "写入响应失败")
		}
		l.Infof("成功代理 m3u8: %s", targetURL)
		return nil
	}

	// 如果是 ts 等媒体文件，直接透传
	// 需要排除的头部（避免重复或冲突）
	excludedHeaders := map[string]bool{
		"content-length":                   true, // 由 io.Copy 自动设置
		"access-control-allow-origin":      true, // 避免重复设置 CORS
		"access-control-allow-methods":     true,
		"access-control-allow-headers":     true,
		"access-control-expose-headers":    true,
		"access-control-allow-credentials": true,
		"transfer-encoding":                true, // 避免分块传输冲突
	}

	// 先删除可能存在的 CORS 头（确保干净）
	w.Header().Del("Access-Control-Allow-Origin")
	w.Header().Del("Access-Control-Allow-Methods")
	w.Header().Del("Access-Control-Allow-Headers")
	w.Header().Del("Access-Control-Expose-Headers")
	w.Header().Del("Access-Control-Allow-Credentials")

	for k, v := range resp.Header {
		lowerKey := strings.ToLower(k)
		if !excludedHeaders[lowerKey] {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	}

	// 设置我们自己的 CORS 头（在复制响应头之后，确保覆盖）
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		l.Errorf("复制响应体失败: %v", err)
		return errs.New(errs.CodeInternalError, "复制响应体失败")
	}

	// 日志记录不同类型的文件
	if strings.HasSuffix(strings.ToLower(targetURL), ".ts") {
		l.Infof("代理 ts 分片: %s", targetURL)
	} else {
		l.Infof("代理资源: %s", targetURL)
	}

	return nil
}

// getDomain 获取请求的域名（处理代理头）
func (l *M3u8ProxyLogic) getDomain(r *http.Request) string {
	// 协议
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	// Host
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = r.Host
	}

	// Nginx 传来的前缀
	prefix := r.Header.Get("X-Forwarded-Prefix")

	return fmt.Sprintf("%s://%s%s", scheme, host, prefix)
}
