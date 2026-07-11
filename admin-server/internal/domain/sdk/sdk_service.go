package sdk

import (
	"context"

	sdkmodel "postapocgame/admin-server/internal/model/sdk"
	"postapocgame/admin-server/internal/repository"
	sdkrepo "postapocgame/admin-server/internal/repository/sdk"
)

type SDKService struct {
	repo *repository.Repository
}

func NewSDKService(repo *repository.Repository) *SDKService {
	return &SDKService{repo: repo}
}

// SaveApiKeyBindings 把"软删旧绑定 + 插入新绑定"包进事务，
// SaveBindings 方法本身不用改，只需要把它构造在事务绑定过的 Repository 之上。
func (s *SDKService) SaveApiKeyBindings(ctx context.Context, sdkKeyID uint64, bindings []sdkmodel.SdkKeyApi) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		return sdkrepo.NewSdkAdminRepository(txRepo).SaveBindings(ctx, sdkKeyID, bindings)
	})
}
