package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DSN string
	}
	// ContentRedis 不叫 Redis：zrpc.RpcServerConf 内嵌字段本身就有一个 Redis（用于可选的
	// gRPC 鉴权），撞名会导致 go-zero 的 conf 解析绑定到错误的结构体，和
	// services/task、services/sdk、services/chat 的 TaskRedis/SdkRedis/ChatRedis 同一个坑。
	// content-rpc 业务本身不用 Redis，纯粹是满足 goctl 生成的 Model 内部 CachedConn 强制要求
	// 非空缓存节点（cache.New 对空 CacheConf 会 log.Fatal），与 gateway 共享同一个 Redis
	// 实例（缓存/锁/队列不拆分，见 16-rpc-conventions.md 第 6 节）。
	ContentRedis struct {
		Address  string
		Password string
	}
	// IamCallbackRpc 连到单体内嵌 pkg/iamcallback.IamCallback server 的 zrpc client 配置。
	// iam 域还没拆分成独立服务前的临时方案，见 pkg/iamcallback 包注释、
	// internal/rpcserver/iamcallback/server.go。content-rpc 用它读取
	// PublicBlogAuthorInfo 需要的用户信息 + 回调写审计日志（RecordAuditLog）。
	IamCallbackRpc zrpc.RpcClientConf
	// Limits 取代原来读字典 blog_*_max_length / blog_article_top_max_count（物理属于 iam 域）
	// 的做法，见 18-service-extraction-runbook.md 2.4 节。默认值全部对齐现有字典种子数据。
	Limits struct {
		BlogTagNameMaxLength          int64 `json:",default=10"`
		BlogArticleTitleMaxLength     int64 `json:",default=100"`
		BlogArticleSummaryLength      int64 `json:",default=120"`
		BlogArticleTopMaxCount        int64 `json:",default=1"`
		BlogFriendLinkNameMaxLength   int64 `json:",default=15"`
		BlogFriendLinkUrlMaxLength    int64 `json:",default=255"`
		BlogFriendLinkRemarkMaxLength int64 `json:",default=127"`
		BlogSocialInfoNameMaxLength   int64 `json:",default=15"`
		BlogSocialInfoUrlMaxLength    int64 `json:",default=255"`
		BlogSocialInfoRemarkMaxLength int64 `json:",default=127"`
	}
}
