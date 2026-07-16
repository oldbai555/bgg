package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/zeromicro/go-zero/zrpc"

	"postapocgame/admin-server/services/iam/iamclient"
)

// 一个简单的种子工具：通过 iam-rpc 的 UserCreate 接口创建默认管理员账号。
// admin_user 表物理属于 iam-rpc，本工具不再直连数据库，而是作为一个普通 gRPC 客户端调用
// iam-rpc（密码 bcrypt 哈希、用户名去重都在 UserService.CreateUser 里完成）。
// 使用方式（在 admin-server 目录）：go run ./cmd/adminseed -endpoint 127.0.0.1:8081 -username admin -password 123456

var (
	endpoint = flag.String("endpoint", "127.0.0.1:8081", "iam-rpc listen address")
	username = flag.String("username", "oldbai", "admin username")
	password = flag.String("password", "oldbai", "admin password (will be bcrypt hashed)")
)

func main() {
	flag.Parse()

	client, err := zrpc.NewClient(zrpc.RpcClientConf{Endpoints: []string{*endpoint}})
	if err != nil {
		log.Fatalf("dial iam-rpc failed: %v", err)
	}
	iamRPC := iamclient.NewIam(client)

	ctx := context.Background()
	_, err = iamRPC.UserCreate(ctx, &iamclient.UserCreateRequest{
		Username: *username,
		Password: *password,
		Status:   1,
	})
	if err != nil {
		if strings.Contains(err.Error(), "已存在") {
			fmt.Printf("Admin user already exists: username=%s\n", *username)
			return
		}
		log.Fatalf("create admin user failed: %v", err)
	}
	fmt.Printf("Admin user created: username=%s\n", *username)
}
