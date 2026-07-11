package errs

import (
	stderr "errors"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error 统一业务错误结构，包含错误码与对外消息。
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// New 创建不带底层堆栈的业务错误。
func New(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// Wrap 为已有错误增加业务码与上下文堆栈。
func Wrap(code int, msg string, err error) *Error {
	if err == nil {
		return New(code, msg)
	}
	return &Error{
		Code:    code,
		Message: msg,
		Err:     errors.Wrap(err, msg),
	}
}

// FromError 尝试从任意 error 中解析出业务错误。
func FromError(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var e *Error
	if stderr.As(err, &e) {
		return e, true
	}
	return nil, false
}

// WrapGRPCError 把 zrpc client 调用返回的 gRPC error 映射成统一的业务 Error，供 gateway
// 侧薄胶水 Logic（如 internal/logic/task/**）在调 RPC 服务后统一处理错误，不用每个调用点
// 各自判断 gRPC status code。未识别的 code 一律归为 CodeInternalError。
func WrapGRPCError(msg string, err error) *Error {
	if err == nil {
		return nil
	}

	code := CodeInternalError
	switch status.Code(err) {
	case codes.NotFound:
		code = CodeNotFound
	case codes.PermissionDenied:
		code = CodeForbidden
	case codes.Unauthenticated:
		code = CodeUnauthorized
	case codes.InvalidArgument, codes.FailedPrecondition:
		code = CodeBadRequest
	case codes.AlreadyExists:
		code = CodeConflict
	case codes.ResourceExhausted:
		code = CodeTooManyRequests
	case codes.Unavailable:
		code = CodeBadGateway
	}
	return Wrap(code, msg, err)
}
