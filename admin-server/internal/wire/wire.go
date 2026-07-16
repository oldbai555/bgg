//go:build wireinject
// +build wireinject

package wire

import (
	"postapocgame/admin-server/internal/config"
	"postapocgame/admin-server/internal/svc"

	"github.com/google/wire"
)

//go:generate go run github.com/google/wire/cmd/wire@v0.6.0

// InitializeApp 编译期注入组合根依赖，返回 ServiceContext 与优雅关闭函数。
func InitializeApp(c config.Config) (*svc.ServiceContext, func(), error) {
	panic(wire.Build(ProviderSet))
}
