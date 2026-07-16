// Package sdk 从 internal/domain/sdk/sdk_service.go 原样搬迁而来，唯一的结构性改动是
// 依赖从单体的 *repository.Repository 换成 sdk-rpc 自己的 *repository.Store。
package sdk

import (
	"context"

	sdkmodel "postapocgame/admin-server/services/sdk/internal/model/sdk"
	"postapocgame/admin-server/services/sdk/internal/repository"
	sdkrepo "postapocgame/admin-server/services/sdk/internal/repository/sdk"
)

type SDKService struct {
	store *repository.Store
}

func NewSDKService(store *repository.Store) *SDKService {
	return &SDKService{store: store}
}

// SaveApiKeyBindings 把"软删旧绑定 + 插入新绑定"包进事务，
// SaveBindings 方法本身不用改，只需要把它构造在事务绑定过的 Store 之上。
func (s *SDKService) SaveApiKeyBindings(ctx context.Context, sdkKeyID uint64, bindings []sdkmodel.SdkKeyApi) error {
	return s.store.Transact(ctx, func(ctx context.Context, txStore *repository.Store) error {
		return sdkrepo.NewSdkAdminRepository(txStore).SaveBindings(ctx, sdkKeyID, bindings)
	})
}
