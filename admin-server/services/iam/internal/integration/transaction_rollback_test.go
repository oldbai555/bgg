//go:build integration

package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	iammodel "postapocgame/admin-server/services/iam/internal/model/iam"
	"postapocgame/admin-server/services/iam/internal/repository"
	iamrepo "postapocgame/admin-server/services/iam/internal/repository/iam"
)

// 08-testing-strategy.md §5 场景 6（可选）：Repository.Transact 回滚在真实 MySQL 上确实生效。
// 单元测试（services/iam/internal/repository/repository_test.go）用 sqlmock 验证了"调用了
// Rollback"，这里补一个真实 MySQL 上的对照：数据真的没有落库，两者互补、不重复。
func TestIntegration_Transact_RollbackOnRealMySQL(t *testing.T) {
	env := setupTestEnv(t)
	ctx := context.Background()

	code := "it_txn_rollback_" + uniqueSuffix()

	txErr := env.Repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		roleRepo := iamrepo.NewRoleRepository(txRepo)
		if err := roleRepo.Create(ctx, &iammodel.AdminRole{
			Name: "事务回滚集成测试角色", Code: code, Status: 1,
		}); err != nil {
			return err
		}
		return errors.New("模拟事务中途失败，触发回滚")
	})
	require.Error(t, txErr)

	roleRepo := iamrepo.NewRoleRepository(env.Repo)
	role, err := roleRepo.FindByCode(ctx, code)
	assert.Error(t, err, "回滚后这条角色记录不应该真的落库")
	assert.Nil(t, role)
}
