// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"time"

	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
)

type SDKCallLogMiddleware struct {
	svcCtx *svc.ServiceContext
}

func NewSDKCallLogMiddleware(svcCtx *svc.ServiceContext) *SDKCallLogMiddleware {
	return &SDKCallLogMiddleware{svcCtx: svcCtx}
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

		logRepo := repository.NewSdkRepository(m.svcCtx.Repository)
		_ = logRepo.SaveCallLog(ctx, &model.SdkCallLog{
			SdkKeyId:       sdkKeyId,
			SdkInterfaceId: sdkInterfaceId,
			ApiCode:        apiCode,
			Path:           r.URL.Path,
			Method:         r.Method,
			Ip:             clientIP,
			UserAgent:      r.UserAgent(),
			ReqBody:        nullString(reqBuf.buf.String()),
			RespBody:       nullString(string(respBody)),
			RespCode:       int64(rec.statusCode),
			DurationMs:     int64(duration / time.Millisecond),
			CreatedAt:      time.Now().Unix(),
			UpdatedAt:      time.Now().Unix(),
		})
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

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
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
