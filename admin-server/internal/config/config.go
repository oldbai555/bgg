// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
//
//staticcheck:ignore SA5008 // "optional" is a go-zero framework extension for JSON tags

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

// Config 聚合服务配置，RestConf 内嵌以支持 go-zero HTTP 配置。
// iam 域拆分成独立服务后，gateway 不再直连任何 MySQL（DatabaseConf 整个类型已删除），
// 只保留共享 Redis（token 黑名单、限流滑动窗口）+ 5 个 zrpc client。
type Config struct {
	rest.RestConf `json:",inline" yaml:",inline" mapstructure:",squash"`
	Redis         RedisConf     `json:"redis,optional" yaml:"redis" mapstructure:"redis"`
	JWT           JWTConf       `json:"jwt,optional" yaml:"jwt" mapstructure:"jwt"`
	Bcrypt        BcryptConf    `json:"bcrypt,optional" yaml:"bcrypt" mapstructure:"bcrypt"`
	RateLimit     RateLimitConf `json:"rateLimit,optional" yaml:"rateLimit" mapstructure:"rateLimit"`
	// IamRpc 连到 iam-rpc（services/iam/）的 zrpc client 配置。iam+system+monitoring+misc
	// 四个域已拆分成独立服务，gateway 侧全部相关 logic、5 个中间件（Auth/Permission/
	// ApiEnabled/OperationLog/Performance）都通过这个 client 调用。
	IamRpc zrpc.RpcClientConf `json:"iamRpc,optional" yaml:"iamRpc" mapstructure:"iamRpc"`
	// IamCallbackRpc 连到 iam-rpc 同一进程内注册的 pkg/iamcallback.IamCallback 服务
	// （RBAC 变更审计日志写入，pkg/audit.RecordAuditLog 用它）。和 IamRpc 指向同一个
	// iam-rpc 地址，只是调用不同的 gRPC service。
	IamCallbackRpc zrpc.RpcClientConf `json:"iamCallbackRpc,optional" yaml:"iamCallbackRpc" mapstructure:"iamCallbackRpc"`
	// TaskRPCConf 连到 task-rpc（services/task/）的 zrpc client 配置。task 域已拆分成独立
	// 服务，gateway 侧的 TaskList/TaskDetail/TaskCancel/TaskRecent 和 5 个导出 logic 都通过
	// 这个 client 调用，不再直接持有 Domain.Task。见 16-rpc-conventions.md 第 5 节。
	TaskRPCConf zrpc.RpcClientConf `json:"taskRpc,optional" yaml:"taskRpc" mapstructure:"taskRpc"`
	// SdkRPCConf 连到 sdk-rpc（services/sdk/）的 zrpc client 配置。sdk 域已拆分成独立服务，
	// gateway 侧的 11 个 SdkApiKey/SdkInterface/SdkCallLog logic 和 SDKAuthMiddleware/
	// SDKRateLimitMiddleware/SDKCallLogMiddleware 三个中间件都通过这个 client 调用，不再
	// 直接持有 Domain.SDK。见 18-service-extraction-runbook.md 2.2 节。
	SdkRPCConf zrpc.RpcClientConf `json:"sdkRpc,optional" yaml:"sdkRpc" mapstructure:"sdkRpc"`
	// ChatRPCConf 连到 chat-rpc（services/chat/）的 zrpc client 配置。chat 域已拆分成独立
	// 服务，gateway 侧的 11 个 Chat*/ChatGroup*/ChatMessage* logic 和 WS↔gRPC 桥接 handler
	// （internal/handler/chat/chatwshandler.go）都通过这个 client 调用，不再直接持有
	// Domain.Chat/ChatHub。见 16-rpc-conventions.md 第 7 节、18-service-extraction-runbook.md
	// 2.3 节。
	ChatRPCConf zrpc.RpcClientConf `json:"chatRpc,optional" yaml:"chatRpc" mapstructure:"chatRpc"`
	// ContentRPCConf 连到 content-rpc（services/content/）的 zrpc client 配置。blog+video 域
	// 已拆分成独立服务，gateway 侧的 34 个 Blog*/PublicBlog* logic + 6 个 Video*/PublicVideo*
	// logic 都通过这个 client 调用，不再直接持有 Domain.Blog/Domain.Video。见
	// 18-service-extraction-runbook.md 2.4 节。
	ContentRPCConf zrpc.RpcClientConf `json:"contentRpc,optional" yaml:"contentRpc" mapstructure:"contentRpc"`
}

type RedisConf struct {
	Address     string `json:"address" yaml:"address" mapstructure:"address"`
	Password    string `json:"password" yaml:"password" mapstructure:"password"`
	DB          int    `json:"db" yaml:"db" mapstructure:"db"`
	Timeout     int    `json:"timeout" yaml:"timeout" mapstructure:"timeout"`             // 连接超时（秒），默认 5
	DialTimeout int    `json:"dialTimeout" yaml:"dialTimeout" mapstructure:"dialTimeout"` // 拨号超时（秒），默认 5
}

type JWTConf struct {
	AccessSecret  string `json:"accessSecret" yaml:"accessSecret" mapstructure:"accessSecret"`
	RefreshSecret string `json:"refreshSecret" yaml:"refreshSecret" mapstructure:"refreshSecret"`
	AccessExpire  int64  `json:"accessExpire" yaml:"accessExpire" mapstructure:"accessExpire"`
	RefreshExpire int64  `json:"refreshExpire" yaml:"refreshExpire" mapstructure:"refreshExpire"`
	Issuer        string `json:"issuer" yaml:"issuer" mapstructure:"issuer"`
}

type BcryptConf struct {
	Cost int `json:"cost" yaml:"cost" mapstructure:"cost"`
}

type RateLimitConf struct {
	Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	IPLimit struct {
		Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
		Quota   int  `json:"quota" yaml:"quota" mapstructure:"quota"`    // 时间窗口内的请求数
		Period  int  `json:"period" yaml:"period" mapstructure:"period"` // 时间窗口（秒）
	} `json:"ipLimit" yaml:"ipLimit" mapstructure:"ipLimit"`
	UserLimit struct {
		Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
		Quota   int  `json:"quota" yaml:"quota" mapstructure:"quota"`    // 时间窗口内的请求数
		Period  int  `json:"period" yaml:"period" mapstructure:"period"` // 时间窗口（秒）
	} `json:"userLimit" yaml:"userLimit" mapstructure:"userLimit"`
	APILimit struct {
		Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
		Quota   int  `json:"quota" yaml:"quota" mapstructure:"quota"`    // 时间窗口内的请求数
		Period  int  `json:"period" yaml:"period" mapstructure:"period"` // 时间窗口（秒）
	} `json:"apiLimit" yaml:"apiLimit" mapstructure:"apiLimit"`
	GlobalLimit struct {
		Enabled bool `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
		Quota   int  `json:"quota" yaml:"quota" mapstructure:"quota"`    // 时间窗口内的请求数
		Period  int  `json:"period" yaml:"period" mapstructure:"period"` // 时间窗口（秒）
	} `json:"globalLimit" yaml:"globalLimit" mapstructure:"globalLimit"`
}
