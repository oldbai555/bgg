package chat

import (
	"context"
	"time"

	"postapocgame/admin-server/internal/consts"
	chatmodel "postapocgame/admin-server/internal/model/chat"
	"postapocgame/admin-server/internal/repository"
	chatrepo "postapocgame/admin-server/internal/repository/chat"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

// Onboarding 供其它领域（当前只有 IAM）触发的窄接口：
// 只暴露"新用户上线要做什么"，不暴露 Chat 域仓储/模型细节，IAM 不需要 import internal/repository/chat。
//
//go:generate mockery --name=Onboarding --output=../../mocks/chat --outpkg=chat_mocks
type Onboarding interface {
	InitNewUser(ctx context.Context, newUserID uint64) error
}

// UserRef 是 Chat 域视角下对"一个用户"的最小引用，不依赖 iammodel.AdminUser，
// 这样 chatdomain 包不需要 import internal/model/iam。
type UserRef struct {
	ID uint64
}

// UserLister 是 Chat 域向外部要求提供的窄接口：分批遍历"活跃、未删除"的用户引用，
// 用于批量建私聊。实现方（IAM 的 UserRepository.FindChunk）在 registry.NewDomain 里适配注入，
// chatdomain 包本身不知道实现方是谁。
//
//go:generate mockery --name=UserLister --output=../../mocks/chat --outpkg=chat_mocks
type UserLister interface {
	FindChunk(ctx context.Context, limit int, lastID uint64) ([]UserRef, uint64, error)
}

type ChatOnboardingService struct {
	repo       *repository.Repository
	userLister UserLister
}

func NewChatOnboardingService(repo *repository.Repository, userLister UserLister) *ChatOnboardingService {
	return &ChatOnboardingService{repo: repo, userLister: userLister}
}

func (s *ChatOnboardingService) InitNewUser(ctx context.Context, newUserID uint64) error {
	s.joinDefaultGroup(ctx, newUserID)
	return s.createPrivateChatsForExistingUsers(ctx, newUserID)
}

// joinDefaultGroup 加入默认企业群组；失败只记日志，不影响后续私聊初始化。
func (s *ChatOnboardingService) joinDefaultGroup(ctx context.Context, newUserID uint64) {
	chatRepo := chatrepo.NewChatRepository(s.repo)
	chatUserRepo := chatrepo.NewChatUserRepository(s.repo)

	groupChat, err := chatRepo.FindByID(ctx, consts.DefaultGroupChatID)
	if err != nil || groupChat.DeletedAt != 0 {
		logx.Infof("默认企业群组不存在或已删除，跳过加入群组操作: userId=%d", newUserID)
		return
	}

	chatUsers, _ := chatUserRepo.FindByChatID(ctx, consts.DefaultGroupChatID)
	for _, cu := range chatUsers {
		if cu.UserId == newUserID {
			return // 已在群组中
		}
	}

	now := time.Now().Unix()
	if err := chatUserRepo.Create(ctx, &chatmodel.ChatUser{
		ChatId: consts.DefaultGroupChatID, UserId: newUserID, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		logx.Errorf("将新用户加入默认企业群组失败: userId=%d, err=%v", newUserID, err)
	}
}

// createPrivateChatsForExistingUsers 用 s.userLister.FindChunk 分批遍历存量用户，
// 每个用户的"建私聊 + 拉两人入会"包一层 s.repo.Transact，批内一个用户失败不影响其他用户，
// 整个方法失败也不影响 IAM 那边已经提交的用户创建（由调用方 UserDomainService 决定"失败只记日志"）。
func (s *ChatOnboardingService) createPrivateChatsForExistingUsers(ctx context.Context, newUserID uint64) error {
	chatRepo := chatrepo.NewChatRepository(s.repo)

	const limit = 100
	lastID := uint64(0)
	for {
		users, newLastID, err := s.userLister.FindChunk(ctx, limit, lastID)
		if err != nil {
			return errs.Wrap(errs.CodeInternalError, "分批查询用户失败", err)
		}
		if len(users) == 0 {
			break
		}

		for _, existingUser := range users {
			if existingUser.ID == newUserID {
				continue
			}
			if existing, err := chatRepo.FindPrivateChatByUserIDs(ctx, newUserID, existingUser.ID); err == nil && existing != nil {
				continue // 私聊已存在
			}
			if err := s.createPrivateChat(ctx, newUserID, existingUser.ID); err != nil {
				logx.Errorf("创建私聊失败: newUserId=%d, existingUserId=%d, err=%v", newUserID, existingUser.ID, err)
				continue
			}
		}

		if len(users) < limit {
			break
		}
		lastID = newLastID
	}
	return nil
}

// createPrivateChat 建一条私聊 + 两条 chat_user 关联，三条写在一个事务里，
// 修复原来"私聊建了、一方或两方没加进去"的部分失败风险。
func (s *ChatOnboardingService) createPrivateChat(ctx context.Context, userA, userB uint64) error {
	return s.repo.Transact(ctx, func(ctx context.Context, txRepo *repository.Repository) error {
		chatRepo := chatrepo.NewChatRepository(txRepo)
		chatUserRepo := chatrepo.NewChatUserRepository(txRepo)

		now := time.Now().Unix()
		privateChat := &chatmodel.Chat{Type: consts.ChatTypePrivate, CreatedBy: 0, CreatedAt: now, UpdatedAt: now}
		if err := chatRepo.Create(ctx, privateChat); err != nil {
			return err
		}
		if err := chatUserRepo.Create(ctx, &chatmodel.ChatUser{
			ChatId: privateChat.Id, UserId: userA, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			return err
		}
		return chatUserRepo.Create(ctx, &chatmodel.ChatUser{
			ChatId: privateChat.Id, UserId: userB, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
		})
	})
}
