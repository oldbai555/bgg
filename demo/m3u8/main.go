package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// CORS 中间件：统一处理跨域请求
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")                           // 允许所有域名跨域访问
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")               // 允许的方法
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization") // 允许的头部

		// 如果是 OPTIONS 请求，直接返回 200 OK
		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusOK)
			return
		}

		c.Next()
	}
}

// 代理处理函数
func proxyHandler(c *gin.Context) {
	// 获取代理的目标 URL
	targetURL := c.DefaultQuery("url", "")
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少url参数"})
		return
	}

	domain := c.DefaultQuery("domain", "http://localhost:8888")

	// 向目标地址发送请求
	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("❌ 请求目标失败: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "请求目标地址失败"})
		return
	}
	defer resp.Body.Close()

	// 如果是 m3u8 文件
	if strings.HasSuffix(strings.ToLower(targetURL), ".m3u8") {
		c.Header("Content-Type", "application/vnd.apple.mpegurl")

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("❌ 读取 m3u8 失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "读取m3u8失败"})
			return
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
				lines[i] = domain + "/proxy?url=" + url.QueryEscape(line)
			}
		}

		newM3U8 := strings.Join(lines, "\n")
		c.Data(http.StatusOK, "application/vnd.apple.mpegurl", []byte(newM3U8))
		log.Printf("✅ 成功代理 m3u8: %s", targetURL)
		return
	}

	// 如果是 ts 等媒体文件，直接透传
	for k, v := range resp.Header {
		for _, vv := range v {
			c.Header(k, vv)
		}
	}
	c.Status(resp.StatusCode)
	_, _ = io.Copy(c.Writer, resp.Body)

	// 日志记录不同类型的文件
	if strings.HasSuffix(strings.ToLower(targetURL), ".ts") {
		log.Printf("📦 代理 ts 分片: %s", targetURL)
	} else {
		log.Printf("🔄 代理资源: %s", targetURL)
	}
}

// 健康检查接口
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func main() {
	// 创建一个 gin 路由
	r := gin.Default()

	// 应用 CORS 中间件
	r.Use(CORS())

	// 注册路由
	r.GET("/proxy", proxyHandler)
	r.GET("/health", healthHandler)

	// 启动服务
	log.Println("🚀 代理服务器启动: http://localhost:8888")
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("❌ 启动失败: %v", err)
	}
}
