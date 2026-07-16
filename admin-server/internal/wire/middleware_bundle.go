package wire

import "github.com/zeromicro/go-zero/rest"

// MiddlewareBundle 聚合 HTTP 中间件 Handle，供 ServiceContext 回填。
type MiddlewareBundle struct {
	Auth                 rest.Middleware
	ApiEnabled           rest.Middleware
	Permission           rest.Middleware
	OperationLog         rest.Middleware
	PublicOperationLog   rest.Middleware
	RateLimit            rest.Middleware
	Performance          rest.Middleware
	Cors                 rest.Middleware
	SDKAuth              rest.Middleware
	SDKRateLimit         rest.Middleware
	SDKCallLog           rest.Middleware
}
