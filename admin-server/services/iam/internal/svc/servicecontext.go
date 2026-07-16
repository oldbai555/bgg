package svc

import (
	"log"

	"postapocgame/admin-server/services/chat/chatclient"
	"postapocgame/admin-server/services/iam/internal/config"
	"postapocgame/admin-server/services/iam/internal/repository"
	"postapocgame/admin-server/services/iam/internal/repository/registry"
	"postapocgame/admin-server/services/sdk/sdkclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	Repository *repository.Repository
	Domain     *registry.Domain
	// SdkRPC 供 TaskCallback server 的 fetchSdkCallLog 分支回调 sdk-rpc.SdkCallLogExport
	// （见 internal/rpcserver/taskcallback/server.go 搬迁前的实现，契约不变）。
	SdkRPC sdkclient.Sdk
	// ChatRPC 供 task 通知消费者（internal/consumer, 原样搬迁）推送 WS 通知。
	ChatRPC chatclient.Chat
}

func NewServiceContext(c config.Config) *ServiceContext {
	repo, err := repository.BuildSources(c)
	if err != nil {
		log.Fatalf("iam-rpc: 初始化数据源失败: %v", err)
	}

	return &ServiceContext{
		Config:     c,
		Repository: repo,
		Domain:     registry.NewDomain(repo),
		SdkRPC:     sdkclient.NewSdk(zrpc.MustNewClient(c.SdkRpc)),
		ChatRPC:    chatclient.NewChat(zrpc.MustNewClient(c.ChatRpc)),
	}
}
