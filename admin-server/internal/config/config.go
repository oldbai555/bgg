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
type Config struct {
	rest.RestConf `json:",inline" yaml:",inline" mapstructure:",squash"`
	Database      DatabaseConf  `json:"database,optional" yaml:"database" mapstructure:"database"`
	Redis         RedisConf     `json:"redis,optional" yaml:"redis" mapstructure:"redis"`
	JWT           JWTConf       `json:"jwt,optional" yaml:"jwt" mapstructure:"jwt"`
	Bcrypt        BcryptConf    `json:"bcrypt,optional" yaml:"bcrypt" mapstructure:"bcrypt"`
	RateLimit     RateLimitConf `json:"rateLimit,optional" yaml:"rateLimit" mapstructure:"rateLimit"`
	// TaskCallbackRPCConf 单体内嵌的 pkg/taskcallback.TaskCallback zrpc server 监听配置，
	// 供 services/task/（task-rpc）回调取导出数据/登记导出文件。见 16-rpc-conventions.md、
	// 17-async-eventing.md 第 1 节。Phase 2 iam-rpc/sdk-rpc 真正拆分后这个 server 原样搬过去。
	TaskCallbackRPCConf zrpc.RpcServerConf `json:"taskCallbackRpc,optional" yaml:"taskCallbackRpc" mapstructure:"taskCallbackRpc"`
	// TaskRPCConf 连到 task-rpc（services/task/）的 zrpc client 配置。task 域已拆分成独立
	// 服务，gateway 侧的 TaskList/TaskDetail/TaskCancel/TaskRecent 和 5 个导出 logic 都通过
	// 这个 client 调用，不再直接持有 Domain.Task。见 16-rpc-conventions.md 第 5 节。
	TaskRPCConf zrpc.RpcClientConf `json:"taskRpc,optional" yaml:"taskRpc" mapstructure:"taskRpc"`
}

type DatabaseConf struct {
	DSN                string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
	MaxOpen            int    `json:"maxOpen" yaml:"maxOpen" mapstructure:"maxOpen"`                                  // 最大打开连接数
	MaxIdle            int    `json:"maxIdle" yaml:"maxIdle" mapstructure:"maxIdle"`                                  // 最大空闲连接数
	ConnMaxLifetime    int    `json:"connMaxLifetime" yaml:"connMaxLifetime" mapstructure:"connMaxLifetime"`          // 连接最大生存时间（秒），默认 300
	ConnMaxIdleTime    int    `json:"connMaxIdleTime" yaml:"connMaxIdleTime" mapstructure:"connMaxIdleTime"`          // 连接最大空闲时间（秒），默认 600
	SlowQueryThreshold int    `json:"slowQueryThreshold" yaml:"slowQueryThreshold" mapstructure:"slowQueryThreshold"` // 慢查询阈值（毫秒），默认 1000
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
