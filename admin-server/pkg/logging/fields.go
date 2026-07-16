package logging

// 标准字段集：所有服务的结构化日志最少要包含这几个字段，
// 便于跨服务用同一个 trace_id 检索、按 service/user_id 过滤。
// trace_id/span_id/service 由 go-zero 的 logx.WithContext(ctx) 和 Setup(serviceName)
// 自动注入，不需要业务代码手动拼接；user_id 需要在能拿到用户身份的地方显式附加。
const (
	FieldTraceID = "trace_id"
	FieldSpanID  = "span_id"
	FieldService = "service"
	FieldUserID  = "user_id"
)
