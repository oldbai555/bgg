// Package chat 从 internal/domain/chat/onboarding.go 迁移而来。原来的 Onboarding/UserLister
// 窄接口 + registry.NewDomain 里的 iamUserListerAdapter 适配层已经不需要了——那一套是
// "IAM 域同进程内触发 Chat 域领域服务"的进程内接口耦合，chat-rpc 拆分后 IAM 侧改成发布
// stream:chat.user.created 事件（见 internal/domain/iam/user_service.go），这里改成消费者
// （services/chat/internal/consumer/chat_user_created_consumer.go）触发 InitNewUser，
// 方法体本身（joinDefaultGroup、createPrivateChatsForExistingUsers、createPrivateChat）
// 原样保留，只是"分批枚举存量用户"从进程内 UserLister.FindChunk 换成回调单体内嵌的
// IamCallback.FindActiveUserChunk（见 pkg/iamcallback，18-service-extraction-runbook.md
// 2.3 节"先在单体里实现,等 iam-rpc 真正拆分时把这个 server 实现原样搬过去"的既定模式）。
package chat

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	chatmodel "postapocgame/admin-server/services/chat/internal/model/chat"
	"postapocgame/admin-server/services/chat/internal/repository"
	chatrepo "postapocgame/admin-server/services/chat/internal/repository/chat"
)

type ChatOnboardingService struct {
	store       *repository.Store
	iamCallback iamcallbackpb.IamCallbackClient
}

func NewChatOnboardingService(store *repository.Store, iamCallback iamcallbackpb.IamCallbackClient) *ChatOnboardingService {
	return &ChatOnboardingService{store: store, iamCallback: iamCallback}
}

func (s *ChatOnboardingService) InitNewUser(ctx context.Context, newUserID uint64) error {
	s.joinDefaultGroup(ctx, newUserID)
	return s.createPrivateChatsForExistingUsers(ctx, newUserID)
}

// joinDefaultGroup 加入默认企业群组；失败只记日志，不影响后续私聊初始化。
func (s *ChatOnboardingService) joinDefaultGroup(ctx context.Context, newUserID uint64) {
	chatRepo := chatrepo.NewChatRepository(s.store)
	chatUserRepo := chatrepo.NewChatUserRepository(s.store)

	groupChat, err := chatRepo.FindByID(ctx, chatconsts.DefaultGroupChatID)
	if err != nil || groupChat.DeletedAt != 0 {
		logx.Infof("默认企业群组不存在或已删除，跳过加入群组操作: userId=%d", newUserID)
		return
	}

	chatUsers, _ := chatUserRepo.FindByChatID(ctx, chatconsts.DefaultGroupChatID)
	for _, cu := range chatUsers {
		if cu.UserId == newUserID {
			return // 已在群组中
		}
	}

	now := time.Now().Unix()
	if err := chatUserRepo.Create(ctx, &chatmodel.ChatUser{
		ChatId: chatconsts.DefaultGroupChatID, UserId: newUserID, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		logx.Errorf("将新用户加入默认企业群组失败: userId=%d, err=%v", newUserID, err)
	}
}

// createPrivateChatsForExistingUsers 分批回调 IamCallback.FindActiveUserChunk 遍历存量用户
// （原来是 s.userLister.FindChunk，现在是跨进程 RPC），每个用户的"建私聊 + 拉两人入会"
// 包一层 s.store.Transact，批内一个用户失败不影响其他用户，整个方法失败也不影响 IAM 那边
// 已经提交的用户创建（由生产者一侧决定"失败只记日志"，见 user_service.go）。
func (s *ChatOnboardingService) createPrivateChatsForExistingUsers(ctx context.Context, newUserID uint64) error {
	chatRepo := chatrepo.NewChatRepository(s.store)

	const limit = 100
	lastID := uint64(0)
	for {
		resp, err := s.iamCallback.FindActiveUserChunk(ctx, &iamcallbackpb.FindActiveUserChunkRequest{
			Limit: limit, LastId: lastID,
		})
		if err != nil {
			return errs.Wrap(errs.CodeInternalError, "分批查询用户失败", err)
		}
		if len(resp.Users) == 0 {
			break
		}

		for _, existingUser := range resp.Users {
			if existingUser.Id == newUserID {
				continue
			}
			if existing, err := chatRepo.FindPrivateChatByUserIDs(ctx, newUserID, existingUser.Id); err == nil && existing != nil {
				continue // 私聊已存在
			}
			if err := s.createPrivateChat(ctx, newUserID, existingUser.Id); err != nil {
				logx.Errorf("创建私聊失败: newUserId=%d, existingUserId=%d, err=%v", newUserID, existingUser.Id, err)
				continue
			}
		}

		if len(resp.Users) < limit {
			break
		}
		lastID = resp.NextLastId
	}
	return nil
}

// createPrivateChat 建一条私聊 + 两条 chat_user 关联，三条写在一个事务里。
func (s *ChatOnboardingService) createPrivateChat(ctx context.Context, userA, userB uint64) error {
	return s.store.Transact(ctx, func(ctx context.Context, txStore *repository.Store) error {
		chatRepo := chatrepo.NewChatRepository(txStore)
		chatUserRepo := chatrepo.NewChatUserRepository(txStore)

		now := time.Now().Unix()
		privateChat := &chatmodel.Chat{Type: chatconsts.ChatTypePrivate, CreatedBy: 0, CreatedAt: now, UpdatedAt: now}
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
