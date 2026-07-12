package svc

import (
	"log"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/chat/internal/config"
	chatconsumer "postapocgame/admin-server/services/chat/internal/consumer"
	chatdomain "postapocgame/admin-server/services/chat/internal/domain/chat"
	"postapocgame/admin-server/services/chat/internal/hub"
	"postapocgame/admin-server/services/chat/internal/repository"
	chatrepo "postapocgame/admin-server/services/chat/internal/repository/chat"
)

type ServiceContext struct {
	Config config.Config
	// Store 是聚合三张表 Model 的句柄，Chat/ChatUser/ChatMessage 三个字段是绑定在它上面的
	// repository；需要事务的 logic（如 ChatGroupCreate）直接调 Store.Transact，和
	// services/sdk 的 SDKService 持有 *repository.Store 是同一个模式。
	Store       *repository.Store
	Chat        chatrepo.ChatRepository
	ChatUser    chatrepo.ChatUserRepository
	ChatMessage chatrepo.ChatMessageRepository
	Hub         *hub.ChatHub
	// IamCallback 回调单体内嵌的 pkg/iamcallback.IamCallback server，供需要展示对方用户
	// 信息（用户名/昵称/头像/部门名/角色名）的 logic 使用，见 pkg/iamcallback 包注释。
	IamCallback iamcallbackpb.IamCallbackClient

	// OnboardingConsumer 消费 stream:chat.user.created，见 chat.go 里的 Start()/Stop() 生命周期
	// 挂钩（和 services/task/task.go 的 ctx.Scheduler.Start() 同一个模式）。
	OnboardingConsumer *chatconsumer.ChatUserCreatedConsumer
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.Mysql.DSN == "" {
		log.Fatalf("chat-rpc: Mysql.DSN 未配置")
	}
	conn := sqlx.NewMysql(c.Mysql.DSN)

	cacheConf := cache.CacheConf{
		{
			RedisConf: redis.RedisConf{Host: c.ChatRedis.Address, Pass: c.ChatRedis.Password, Type: "node"},
			Weight:    100,
		},
	}
	store := repository.NewStore(conn, cacheConf)

	iamCallbackClient := iamcallbackpb.NewIamCallbackClient(zrpc.MustNewClient(c.IamCallbackRpc).Conn())
	onboarding := chatdomain.NewChatOnboardingService(store, iamCallbackClient)

	rds := redis.MustNewRedis(redis.RedisConf{Host: c.ChatRedis.Address, Pass: c.ChatRedis.Password, Type: "node"})

	chatHub := hub.NewChatHub()
	go chatHub.Run()

	return &ServiceContext{
		Config:             c,
		Store:              store,
		Chat:               chatrepo.NewChatRepository(store),
		ChatUser:           chatrepo.NewChatUserRepository(store),
		ChatMessage:        chatrepo.NewChatMessageRepository(store),
		Hub:                chatHub,
		IamCallback:        iamCallbackClient,
		OnboardingConsumer: chatconsumer.NewChatUserCreatedConsumer(rds, onboarding),
	}
}
