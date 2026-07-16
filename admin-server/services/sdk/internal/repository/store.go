// Package repository 从 internal/repository/sdk/ 原样搬迁而来。唯一的结构性改动是把两个
// repository（SdkAdminRepository/SdkRepository）原来共享的 *repository.Repository（单体聚合了
// 全部 9 个业务域 Model 的大句柄）换成这里的 *Store——sdk-rpc 从第一天起只有 sdk_key/
// sdk_interface/sdk_key_api/sdk_call_log 四张表，不该也不能继续持有指向其它域的句柄。
package repository

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
)

// Store 聚合 sdk-rpc 自己需要的全部 Model，供 SdkAdminRepository/SdkRepository 共用。
type Store struct {
	DB                sqlx.SqlConn
	SdkKeyModel       sdkmodel.SdkKeyModel
	SdkInterfaceModel sdkmodel.SdkInterfaceModel
	SdkKeyApiModel    sdkmodel.SdkKeyApiModel
	SdkCallLogModel   sdkmodel.SdkCallLogModel
}

func NewStore(conn sqlx.SqlConn, cacheConf cache.CacheConf) *Store {
	return &Store{
		DB:                conn,
		SdkKeyModel:       sdkmodel.NewSdkKeyModel(conn, cacheConf),
		SdkInterfaceModel: sdkmodel.NewSdkInterfaceModel(conn, cacheConf),
		SdkKeyApiModel:    sdkmodel.NewSdkKeyApiModel(conn, cacheConf),
		SdkCallLogModel:   sdkmodel.NewSdkCallLogModel(conn, cacheConf),
	}
}

// Transact 在单个 MySQL 事务内执行 fn，用法和 internal/repository/repository.go 的
// Repository.Transact 完全同构（sdk-rpc 自己的小号版本，只换绑这四个 Model）。
func (s *Store) Transact(ctx context.Context, fn func(ctx context.Context, txStore *Store) error) error {
	return s.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, s.withSession(session))
	})
}

func (s *Store) withSession(session sqlx.Session) *Store {
	return &Store{
		DB:                sqlx.NewSqlConnFromSession(session),
		SdkKeyModel:       s.SdkKeyModel.WithSession(session),
		SdkInterfaceModel: s.SdkInterfaceModel.WithSession(session),
		SdkKeyApiModel:    s.SdkKeyApiModel.WithSession(session),
		SdkCallLogModel:   s.SdkCallLogModel.WithSession(session),
	}
}
