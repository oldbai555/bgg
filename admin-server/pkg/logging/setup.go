package logging

import (
	"github.com/zeromicro/go-zero/core/logx"
)

// Setup 统一的日志初始化，六个服务（gateway + iam/task/sdk/chat/content-rpc）的 main
// 函数都调用这一个函数，不再各自内联 logx.SetUp(logx.LogConf{...})。
// serviceName 写入每条日志的 service 字段，用于在聚合后按服务过滤。
func Setup(serviceName string) error {
	return logx.SetUp(logx.LogConf{
		ServiceName: serviceName,
		Encoding:    "json",
		Level:       "info",
	})
}
