package consts

import "time"

// 通用状态字符串
const (
	StatusOK    = "ok"
	StatusError = "error"
)

// Ping 接口相关常量
const (
	PingMessagePong = "pong"
)

// 字典 code 常量（物理属于 iam-rpc 的 admin_dict_type/admin_dict_item 表，
// gateway 侧通过 IamRPC.DictGet 读取，这里只是本地引用用的 code 常量）
const (
	DictCodeVideoProxyURL  = "video_proxy_url"
	DictCodeStorageBaseURL = "storage_base_url"
)

// 需要 CORS 的公共接口路径片段（ApiEnabledMiddleware.setCORSIfNeeded 用，用
// strings.Contains 匹配；和下面已存在、构造完整代理 URL 用的 PathM3U8Proxy 是不同用途，
// 不合并）
const (
	CORSPathM3U8Proxy    = "/m3u8/proxy"
	CORSPathVideoCollect = "/videos/collect"
)

// Redis 相关常量
const (
	RedisPingFailedMessage = "redis ping failed"

	// JWT 黑名单前缀
	RedisJWTBlacklistPrefix = "jwt:blacklist:"

	// 限流相关 Redis 前缀
	RedisRateLimitGlobalPrefix = "rate_limit:global"
	RedisRateLimitIPPrefix     = "rate_limit:ip:"
	RedisRateLimitUserPrefix   = "rate_limit:user:"
	RedisRateLimitAPIPrefix    = "rate_limit:api:"
)

// 限流提示信息
const (
	RateLimitMessageGlobal = "请求过于频繁，请稍后再试（全局限流）"
	RateLimitMessageIP     = "请求过于频繁，请稍后再试（IP限流）"
	RateLimitMessageUser   = "请求过于频繁，请稍后再试（用户限流）"
	RateLimitMessageAPI    = "请求过于频繁，请稍后再试（接口限流）"
)

// 常用路径常量
const (
	PathPing = "/api/v1/ping"

	// 认证相关路径
	PathLogin   = "/api/v1/login"
	PathLogout  = "/api/v1/logout"
	PathRefresh = "/api/v1/refresh"

	// WebSocket 路径
	PathChatWS = "/api/v1/chats/ws"

	// 文件上传相关路径
	PathFileUploads = "/api/v1/files/uploads"
)

// 文件系统路径常量
const (
	// UploadDir 上传文件存储目录
	UploadDir = "./uploads"
)

// 公告状态常量
const (
	// NoticeStatusDraft 草稿
	NoticeStatusDraft int64 = 1
	// NoticeStatusPublished 已发布
	NoticeStatusPublished int64 = 2
)

const (
	Open = 1
)

// 用户状态常量
const (
	// UserStatusEnabled 启用
	UserStatusEnabled int64 = 1
)

// 任务状态常量（对应字典 task_status 的 value）
const (
	// TaskStatusPending 未开始
	TaskStatusPending int64 = 1
	// TaskStatusRunning 进行中
	TaskStatusRunning int64 = 2
	// TaskStatusCompleted 已完成
	TaskStatusCompleted int64 = 3
	// TaskStatusFailed 失败
	TaskStatusFailed int64 = 4
)

// 任务类型常量（对应字典 task_type 的 value）
const (
	// TaskTypeExcelExport 异步导出Excel
	TaskTypeExcelExport int64 = 1
)

// 任务执行类型常量（对应字典 task_execution_type 的 value）
const (
	// TaskExecutionTypeSync 同步执行
	TaskExecutionTypeSync int64 = 1
	// TaskExecutionTypeAsync 异步执行
	TaskExecutionTypeAsync int64 = 2
)

// 任务通知来源类型常量（对应字典 notification_source_type 的 value）
const (
	// NotificationSourceTypeTask 任务通知来源类型
	NotificationSourceTypeTask = "task"
)

// 任务相关 Redis 键前缀
const (
	// RedisTaskLockPrefix 任务锁前缀（用于分布式锁）
	RedisTaskLockPrefix = "task:lock:"
	// RedisTaskConfigPrefix 任务配置前缀（用于存储任务配置）
	RedisTaskConfigPrefix = "task:config:"
)

// 任务通知消息常量
const (
	// TaskNotificationTitleRunning 任务执行中通知标题
	TaskNotificationTitleRunning = "任务执行中"
	// TaskNotificationTitleCompleted 任务完成通知标题
	TaskNotificationTitleCompleted = "任务执行完成"
	// TaskNotificationTitleFailed 任务失败通知标题
	TaskNotificationTitleFailed = "任务执行失败"
)

// 任务配置常量（默认值）
const (
	// TaskDefaultScanInterval 默认扫描间隔（秒）
	TaskDefaultScanInterval = 5
	// TaskDefaultMaxConcurrent 默认最大并发执行数量
	TaskDefaultMaxConcurrent = 10
	// TaskDefaultBatchSize 默认批次大小
	TaskDefaultBatchSize = 100
	// TaskDefaultTaskTimeout 默认任务超时时间（秒，30分钟）
	TaskDefaultTaskTimeout = 1800
	// TaskDefaultLockTimeout 默认锁超时时间（秒）
	TaskDefaultLockTimeout = 30
)

// 任务相关 WebSocket 消息类型常量
const (
	// WSTaskProgress 任务进度消息类型
	WSTaskProgress = "task_progress"
	// WSNotification 通知消息类型
	WSNotification = "notification"
)

// 任务通知级别常量
const (
	// TaskNotificationLevelInfo 信息级别
	TaskNotificationLevelInfo = "info"
	// TaskNotificationLevelSuccess 成功级别
	TaskNotificationLevelSuccess = "success"
	// TaskNotificationLevelError 错误级别
	TaskNotificationLevelError = "error"
)

// 任务导出模块常量（对应各业务模块的标识）
const (
	// TaskModuleOperationLog 操作日志导出模块
	TaskModuleOperationLog = "operation_log"
	// TaskModuleAuditLog 审计日志导出模块
	TaskModuleAuditLog = "audit_log"
	// TaskModuleLoginLog 登录日志导出模块
	TaskModuleLoginLog = "login_log"
	// TaskModuleSdkCallLog SDK 调用日志导出模块
	TaskModuleSdkCallLog = "sdk_call_log"
	// TaskModulePerformanceLog 性能监控日志导出模块
	TaskModulePerformanceLog = "performance_log"
)

// 任务导出筛选字段 key 常量（ExcelExportParams.Filters 使用）
const (
	TaskFilterUserId        = "userId"
	TaskFilterUsername      = "username"
	TaskFilterOperationType = "operationType"
	TaskFilterOperationObj  = "operationObject"
	TaskFilterMethod        = "method"
	TaskFilterPath          = "path"
	TaskFilterStartTime     = "startTime"
	TaskFilterEndTime       = "endTime"
	TaskFilterAuditType     = "auditType"
	TaskFilterAuditObject   = "auditObject"
	TaskFilterStatus        = "status"
	TaskFilterIsSlow        = "isSlow"
	TaskFilterStatusCode    = "statusCode"
	TaskFilterSdkKeyId      = "sdkKeyId"
	TaskFilterApiCode       = "apiCode"
	TaskFilterRespCode      = "respCode"
	TaskFilterIP            = "ip"
)

// M3U8 代理相关常量
const (
	// M3U8 文件扩展名
	M3U8FileExtension = ".m3u8"

	// HTTP 协议
	ProtocolHTTP  = "http"
	ProtocolHTTPS = "https"

	// HTTP 代理头
	HeaderXForwardedProto  = "X-Forwarded-Proto"
	HeaderXForwardedHost   = "X-Forwarded-Host"
	HeaderXForwardedPrefix = "X-Forwarded-Prefix"

	// HTTP 响应头
	HeaderContentType                   = "Content-Type"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"

	// CORS 值
	CORSAllowOriginAll      = "*"
	CORSAllowMethodsGetOpts = "GET, OPTIONS"
	CORSAllowHeadersDefault = "Content-Type,Authorization"

	// Content-Type
	ContentTypeM3U8        = "application/vnd.apple.mpegurl"
	ContentTypeOctetStream = "application/octet-stream"

	// M3U8 文件格式相关
	M3U8CommentPrefix = "#"
	M3U8LineSeparator = "\n"
	M3U8URLSeparator  = "|"

	// M3U8 代理路径
	PathM3U8Proxy = "/api/v1/m3u8/proxy?url="

	// HTTP 响应头（小写，用于排除）
	HeaderContentLength                   = "content-length"
	HeaderAccessControlAllowOriginLC      = "access-control-allow-origin"
	HeaderAccessControlAllowMethodsLC     = "access-control-allow-methods"
	HeaderAccessControlAllowHeadersLC     = "access-control-allow-headers"
	HeaderAccessControlExposeHeadersLC    = "access-control-expose-headers"
	HeaderAccessControlAllowCredentialsLC = "access-control-allow-credentials"
	HeaderTransferEncoding                = "transfer-encoding"

	// 文件扩展名
	FileExtensionTS = ".ts"
)

// M3U8 代理错误消息
const (
	ErrMsgMissingURL          = "缺少url参数"
	ErrMsgInvalidURLFormat    = "URL格式无效"
	ErrMsgUnsupportedProtocol = "仅支持 http 和 https 协议"
	ErrMsgWriteResponseFailed = "写入响应失败"
	ErrMsgRequestTargetFailed = "请求目标地址失败"
	ErrMsgReadM3U8Failed      = "读取m3u8失败"
	ErrMsgReadResourceFailed  = "读取资源失败"
	ErrMsgCacheNotInitialized = "缓存未初始化"
)

// M3U8 代理日志消息
const (
	LogMsgWriteCacheResponseFailed = "写入缓存响应失败: %v"
	LogMsgReturnFromCache          = "从缓存返回代理内容: %s"
	LogMsgRequestTargetFailed      = "请求目标失败: %v"
	LogMsgReadM3U8Failed           = "读取 m3u8 失败: %v"
	LogMsgCacheM3U8Failed          = "缓存 m3u8 失败: %v"
	LogMsgCacheResourceFailed      = "缓存资源失败: %v"
	LogMsgCacheReadFailed          = "缓存读取失败: %v"
	LogMsgWriteResponseFailed      = "写入响应失败: %v"
	LogMsgProxyM3U8Success         = "成功代理 m3u8: %s"
	LogMsgProxyTSFragment          = "代理 ts 分片: %s"
	LogMsgProxyResource            = "代理资源: %s"
)

// 缓存读取超时时间（500ms，避免 Redis 慢查询阻塞请求）
const CacheReadTimeout = 500 * time.Millisecond

// 代理请求超时时间（10秒，避免请求源服务器时间过长）
const ProxyRequestTimeout = 10 * time.Second

// 代理资源缓存大小限制（5MB，超过此大小的文件不缓存，避免内存和 Redis 压力）
const ProxyCacheMaxSize = 5 * 1024 * 1024 // 5MB
