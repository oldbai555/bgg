package iamcallback

import (
	"github.com/zeromicro/go-zero/zrpc"

	pb "postapocgame/admin-server/pkg/iamcallback/pb"
)

// Client 是 pb.IamCallbackClient 的别名，供 chat-rpc 回调单体内嵌的 IamCallback server
// 取存量用户 ID / 用户展示信息。等 iam-rpc 真正拆分后，Client 构造方式不变，只是连接目标
// 从单体换成 iam-rpc 自己（和 pkg/taskcallback.Client 的处理方式同构）。
type Client = pb.IamCallbackClient

// NewClient 用 zrpc.RpcClientConf 构造一个 IamCallback 客户端。
func NewClient(c zrpc.RpcClientConf) (Client, error) {
	conn, err := zrpc.NewClient(c)
	if err != nil {
		return nil, err
	}
	return pb.NewIamCallbackClient(conn.Conn()), nil
}
