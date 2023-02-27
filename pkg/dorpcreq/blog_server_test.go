package dorpcreq

import (
	"testing"
)

func TestServer(t *testing.T) {
	// 2.实例化gRPC
	// grpc.NewServer(
	// 	grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 		md, b := metadata.FromIncomingContext(ctx)
	// 		if b {
	// 			strings := md.Get("token")
	// 			log.Infof("value: %v", strings)
	// 			return nil, lbuser.ErrUserNotFound
	// 		}
	// 		return handler(ctx, req)
	// 	}))
}

func TestLbblogServer_GetArticleList(t *testing.T) {
	// 1.连接
	// conn, err := grpc.Dial("127.0.0.1:18001", grpc.WithInsecure(), grpc.WithBlock())
	// if err != nil {
	// 	fmt.Printf("连接异常： %s\n", err)
	// }
	// defer conn.Close()
	//
	// // 4. 调用接口
	// // 创建metadata到context中.
	// md := metadata.Pairs("token", time.Now().Format(time.StampNano))
	// ctx := metadata.NewOutgoingContext(context.Background(), md)

}
