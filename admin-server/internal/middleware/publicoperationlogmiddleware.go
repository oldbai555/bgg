package middleware

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// PublicOperationLogMiddleware 公共接口操作日志中间件，记录所有公共接口的调用日志（包括 GET 请求）
type PublicOperationLogMiddleware struct {
	svcCtx *svc.ServiceContext
	logCh  chan *model.AdminOperationLog // 异步日志通道
}

func NewPublicOperationLogMiddleware(svcCtx *svc.ServiceContext) *PublicOperationLogMiddleware {
	m := &PublicOperationLogMiddleware{
		svcCtx: svcCtx,
		logCh:  make(chan *model.AdminOperationLog, 1000), // 缓冲1000条日志
	}

	// 启动异步日志写入 goroutine
	go m.logWriter()

	return m
}

// Handle 中间件处理函数
func (m *PublicOperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 记录所有 HTTP 方法（包括 GET）
		method := r.Method

		// 排除一些不需要记录的接口
		path := r.URL.Path
		if m.shouldSkip(path) {
			next(w, r)
			return
		}

		// 记录开始时间
		startTime := time.Now()

		// 读取请求体（用于记录请求参数）
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // 恢复 body，供后续处理使用
		}

		// 包装 ResponseWriter 以捕获响应
		responseWriter := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           &bytes.Buffer{},
		}

		// 执行下一个处理器
		next(responseWriter, r)

		// 计算耗时
		duration := int(time.Since(startTime).Milliseconds())

		// 公共接口不需要用户认证，userId 和 username 为 0 和空字符串
		userId := uint64(0)
		username := ""

		// 解析操作类型和操作对象
		operationType, operationObject := m.parseOperation(method, path)

		// 构建操作日志
		requestParams := sql.NullString{}
		if len(requestBody) > 0 {
			// 限制请求参数长度，避免过长
			paramsStr := string(requestBody)
			if len(paramsStr) > 10000 {
				paramsStr = paramsStr[:10000] + "..."
			}
			requestParams = sql.NullString{String: paramsStr, Valid: true}
		}

		// 对于 GET 请求，记录查询参数
		if method == http.MethodGet && len(r.URL.RawQuery) > 0 {
			queryStr := r.URL.RawQuery
			if len(queryStr) > 10000 {
				queryStr = queryStr[:10000] + "..."
			}
			if requestParams.Valid {
				requestParams.String = requestParams.String + "&" + queryStr
			} else {
				requestParams = sql.NullString{String: queryStr, Valid: true}
			}
		}

		operationLog := &model.AdminOperationLog{
			UserId:          userId,
			Username:        username,
			OperationType:   operationType,
			OperationObject: operationObject,
			Method:          method,
			Path:            path,
			RequestParams:   requestParams,
			ResponseCode:    int64(responseWriter.statusCode),
			ResponseMsg:     m.extractResponseMsg(responseWriter.body.String()),
			IpAddress:       m.getClientIP(r),
			UserAgent:       r.UserAgent(),
			Duration:        int64(duration),
			DeletedAt:       0, // 软删除字段，0 表示未删除
		}

		// 异步写入日志（非阻塞）
		select {
		case m.logCh <- operationLog:
			logx.Infof("公共接口操作日志已加入队列: method=%s, path=%s", method, path)
		default:
			// 通道满了，记录警告但不阻塞请求
			logx.Errorf("公共接口操作日志通道已满，丢弃日志: %+v", operationLog)
		}
	}
}

// shouldSkip 判断是否应该跳过记录
func (m *PublicOperationLogMiddleware) shouldSkip(path string) bool {
	// 公共接口通常不需要跳过，但可以在这里添加需要跳过的路径
	skipPaths := []string{
		// 可以在这里添加需要跳过的公共接口路径
	}
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// parseOperation 解析操作类型和操作对象
func (m *PublicOperationLogMiddleware) parseOperation(method, path string) (operationType, operationObject string) {
	// 根据 HTTP 方法确定操作类型
	switch method {
	case http.MethodGet:
		operationType = "query" // 公共接口的 GET 请求标记为 query
	case http.MethodPost:
		operationType = "create"
	case http.MethodPut:
		operationType = "update"
	case http.MethodDelete:
		operationType = "delete"
	default:
		operationType = "unknown"
	}

	// 从路径中提取操作对象（模块名）
	// 例如：/api/v1/public/videos/list -> public_video
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 4 {
		// 移除 /api/v1 前缀，获取 public/videos
		module := parts[2] + "_" + parts[3]
		// 移除复数形式（如 videos -> video）
		if strings.HasSuffix(parts[3], "s") && len(parts[3]) > 1 {
			module = parts[2] + "_" + parts[3][:len(parts[3])-1]
		}
		operationObject = module
	} else if len(parts) >= 3 {
		// 如果只有 3 部分，使用第二部分作为操作对象
		module := parts[2]
		// 移除复数形式
		if strings.HasSuffix(module, "s") && len(module) > 1 {
			module = module[:len(module)-1]
		}
		operationObject = module
	}

	return operationType, operationObject
}

// extractResponseMsg 从响应体中提取消息
func (m *PublicOperationLogMiddleware) extractResponseMsg(responseBody string) string {
	if responseBody == "" {
		return ""
	}

	// 尝试解析 JSON 响应
	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &resp); err == nil {
		if msg, ok := resp["msg"].(string); ok {
			return msg
		}
	}

	// 如果解析失败，返回前255个字符
	if len(responseBody) > 255 {
		return responseBody[:255]
	}
	return responseBody
}

// getClientIP 获取客户端 IP 地址
func (m *PublicOperationLogMiddleware) getClientIP(r *http.Request) string {
	// 优先从 X-Forwarded-For 获取（代理场景）
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 其次从 X-Real-IP 获取
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 最后从 RemoteAddr 获取
	ip = r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// logWriter 异步日志写入器
func (m *PublicOperationLogMiddleware) logWriter() {
	operationLogRepo := repository.NewOperationLogRepository(m.svcCtx.Repository)
	batch := make([]*model.AdminOperationLog, 0, 100) // 批量写入，每批100条
	ticker := time.NewTicker(5 * time.Second)         // 每5秒写入一次
	defer ticker.Stop()

	for {
		select {
		case log := <-m.logCh:
			batch = append(batch, log)
			// 如果批次达到100条，立即写入
			if len(batch) >= 100 {
				m.writeBatch(operationLogRepo, batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			// 定时写入批次中的日志
			if len(batch) > 0 {
				m.writeBatch(operationLogRepo, batch)
				batch = batch[:0]
			}
		}
	}
}

// writeBatch 批量写入日志
func (m *PublicOperationLogMiddleware) writeBatch(repo repository.OperationLogRepository, logs []*model.AdminOperationLog) {
	ctx := context.Background()
	if len(logs) == 0 {
		return
	}

	// 使用批量创建方法
	if err := repo.BatchCreate(ctx, logs); err != nil {
		logx.Errorf("批量写入公共接口操作日志失败: count=%d, error: %v", len(logs), err)
		// 如果批量写入失败，尝试逐条写入
		for _, log := range logs {
			if err := repo.Create(ctx, log); err != nil {
				logx.Errorf("写入公共接口操作日志失败: %+v, error: %v", log, err)
			}
		}
	} else {
		logx.Infof("成功批量写入公共接口操作日志: count=%d", len(logs))
	}
}
