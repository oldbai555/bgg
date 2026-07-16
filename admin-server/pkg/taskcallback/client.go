package taskcallback

import (
	"github.com/zeromicro/go-zero/zrpc"

	pb "postapocgame/admin-server/pkg/taskcallback/pb"
)

// Client 是 pb.TaskCallbackClient 的别名，供 services/task/（task-rpc）回调 iam-rpc/sdk-rpc
// 取导出数据/登记导出文件。当前阶段单体内嵌 TaskCallback server（internal/rpcserver/
// taskcallback/），Client 连的就是这个单体进程；Phase 2 iam-rpc/sdk-rpc 真正拆分后，
// task-rpc 侧的 moduleServiceRoute 会指向各自服务，Client 构造方式不变。
type Client = pb.TaskCallbackClient

// NewClient 用 zrpc.RpcClientConf 构造一个 TaskCallback 客户端。
func NewClient(c zrpc.RpcClientConf) (Client, error) {
	conn, err := zrpc.NewClient(c)
	if err != nil {
		return nil, err
	}
	return pb.NewTaskCallbackClient(conn.Conn()), nil
}
