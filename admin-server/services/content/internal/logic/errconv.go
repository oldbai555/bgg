package logic

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"postapocgame/admin-server/pkg/errs"
)

// toGRPCStatus 把内部调用（repository/domain 层）返回的 *errs.Error 转成 gRPC status error。
// 不转换的话，*errs.Error 原样穿过 gRPC 边界会被 gateway 侧 errs.WrapGRPCError
// （pkg/errs/errors.go）识别成未映射的 gRPC code，一律退化成 CodeInternalError——
// 等于白建了这套错误码映射。和 services/task、services/sdk、services/chat 的
// toGRPCStatus 同一份实现（16-rpc-conventions.md 第 6 节"直接复制不共享"），未识别的
// *errs.Error.Code 统一归为 codes.Internal。
func toGRPCStatus(err error) error {
	if err == nil {
		return nil
	}
	e, ok := errs.FromError(err)
	if !ok {
		return err
	}

	code := codes.Internal
	switch e.Code {
	case errs.CodeNotFound:
		code = codes.NotFound
	case errs.CodeBadRequest:
		code = codes.InvalidArgument
	case errs.CodeForbidden:
		code = codes.PermissionDenied
	case errs.CodeUnauthorized:
		code = codes.Unauthenticated
	case errs.CodeConflict:
		code = codes.AlreadyExists
	case errs.CodeTooManyRequests:
		code = codes.ResourceExhausted
	case errs.CodeBadGateway:
		code = codes.Unavailable
	}
	return status.Error(code, e.Message)
}
