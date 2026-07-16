// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/services/sdk/sdkclient"
)

// SDKCallLogMiddleware 继续留在 gateway，内部实现从直连 Repository 改成调 sdk-rpc 的
// RecordCallLog。原代码用 `_ = logRepo.SaveCallLog(...)` 丢弃错误——RPC 边界另一侧
// （services/sdk/internal/logic/recordcallloglogic.go）如实把错误传回来，"失败不影响
// 本次 SDK 调用"这个既有语义在这里（调用方）继续保持：只记日志，不影响响应。
type SDKCallLogMiddleware struct {
	sdkRPC sdkclient.Sdk
}

func NewSDKCallLogMiddleware(sdkRPC sdkclient.Sdk) *SDKCallLogMiddleware {
	return &SDKCallLogMiddleware{sdkRPC: sdkRPC}
}

func (m *SDKCallLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// 通过 TeeReader 收集部分请求体（2KB），不影响后续读取
		var reqBuf limitedBuffer
		reqBuf.limit = 2048
		if r.Body != nil {
			r.Body = io.NopCloser(io.TeeReader(r.Body, &reqBuf))
		}

		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next(rec, r)

		duration := time.Since(start)
		respBody := rec.Body()
		if len(respBody) > 2048 {
			respBody = respBody[:2048]
		}

		sdkKeyId, _ := ctx.Value(ctxKeySdkKeyID).(uint64)
		apiCode, _ := ctx.Value(ctxKeySdkApiCode).(string)
		sdkInterfaceId, _ := ctx.Value(ctxKeySdkInterfaceID).(uint64)
		clientIP := clientIPFromRequest(r)

		_, err := m.sdkRPC.RecordCallLog(ctx, &sdkclient.RecordCallLogRequest{
			SdkKeyId:       sdkKeyId,
			SdkInterfaceId: sdkInterfaceId,
			ApiCode:        apiCode,
			Path:           r.URL.Path,
			Method:         r.Method,
			Ip:             clientIP,
			UserAgent:      r.UserAgent(),
			ReqBody:        reqBuf.buf.String(),
			RespBody:       string(respBody),
			RespCode:       int64(rec.statusCode),
			DurationMs:     int64(duration / time.Millisecond),
		})
		if err != nil {
			logx.Errorf("记录 SDK 调用日志失败: %v", err)
		}
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) Body() []byte {
	return r.body.Bytes()
}

type limitedBuffer struct {
	buf   bytes.Buffer
	limit int
}

func (l *limitedBuffer) Write(p []byte) (int, error) {
	if l.limit <= 0 {
		return len(p), nil
	}
	remain := l.limit - l.buf.Len()
	if remain > 0 {
		if len(p) <= remain {
			_, _ = l.buf.Write(p)
		} else {
			_, _ = l.buf.Write(p[:remain])
		}
	}
	return len(p), nil
}
