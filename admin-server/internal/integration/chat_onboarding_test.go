//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"postapocgame/admin-server/internal/consts"
	iamdomain "postapocgame/admin-server/internal/domain/iam"
)

// 08-testing-strategy.md §5 场景 4：IAM 用户创建 → chat onboarding 全链路。
// 验证 04-domain-iam-chat.md 任务 1 确定的"异步、尽力而为"语义：CreateUser 本身不因
// chat 初始化失败/未完成而失败，且默认企业群组的加入在短暂等待后确实落库。
func TestIntegration_UserCreate_TriggersChatOnboarding(t *testing.T) {
	env := setupTestEnv(t)
	ctx := context.Background()

	username := "it_onboard_" + uniqueSuffix()
	svc := env.Domain.IAM.UserService
	require.NotNil(t, svc, "UserService 需要在 registry.NewDomain 里已经接好 chat onboarding 依赖")

	user, err := svc.CreateUser(ctx, iamdomain.CreateUserInput{
		Username: username,
		Nickname: "Onboarding 集成测试用户",
		Password: "S3cret!Pwd",
		Status:   1,
	})
	require.NoError(t, err)
	require.NotZero(t, user.Id)

	// CreateUser 内部用 go func 异步触发 chatdomain.Onboarding.InitNewUser，这里没有 mockery
	// 可以挂钩的信号，用短暂轮询代替固定 sleep，减少在慢速 CI 环境下的假阴性。
	deadline := time.Now().Add(3 * time.Second)
	var joined bool
	for time.Now().Before(deadline) {
		chatUsers, err := env.Domain.Chat.ChatUser.FindByChatID(ctx, consts.DefaultGroupChatID)
		require.NoError(t, err)
		for _, cu := range chatUsers {
			if cu.UserId == user.Id {
				joined = true
				break
			}
		}
		if joined {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	assert.True(t, joined, "新用户应该在异步 onboarding 完成后出现在默认企业群组成员列表里")
}
