package errs

const (
	// CodeOK 成功。
	CodeOK = 0

	// 通用错误码（1xxxx）
	CodeInternalError   = 10001
	CodeBadRequest      = 10002
	CodeUnauthorized    = 10003
	CodeForbidden       = 10004
	CodeNotFound        = 10005
	CodeBadDB           = 10006
	CodeBadGateway      = 10007 // 网关错误（502）
	CodeConflict        = 10008 // 冲突错误（409）
	CodeTooManyRequests = 10009 // 限流（429）
)
