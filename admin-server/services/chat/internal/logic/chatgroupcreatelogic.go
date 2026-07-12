package logic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/chat/chat"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	chatmodel "postapocgame/admin-server/services/chat/internal/model/chat"
	"postapocgame/admin-server/services/chat/internal/repository"
	chatrepo "postapocgame/admin-server/services/chat/internal/repository/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupCreateLogic {
	return &ChatGroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupCreate 迁移自 internal/logic/chat/group/chat_group_create_logic.go：建群组 + 拉
// 创建人入群两条写包一个事务（Week4-5 修的孤儿群组 bug，见 docs/progress.md），初始成员的
// 存在性校验从本地 Domain.IAM.User.FindByID 改成回调 IamCallback.GetUserProfile。
func (l *ChatGroupCreateLogic) ChatGroupCreate(in *chat.ChatGroupCreateRequest) (*chat.Empty, error) {
	if in.Name == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "群组名称不能为空"))
	}
	if in.OperatorUserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}

	now := time.Now().Unix()
	chatEntity := &chatmodel.Chat{
		Name:        in.Name,
		Type:        chatconsts.ChatTypeGroup,
		Avatar:      in.Avatar,
		Description: in.Description,
		CreatedBy:   in.OperatorUserId,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   0,
	}

	err := l.svcCtx.Store.Transact(l.ctx, func(ctx context.Context, txStore *repository.Store) error {
		if err := chatrepo.NewChatRepository(txStore).Create(ctx, chatEntity); err != nil {
			return errs.Wrap(errs.CodeInternalError, "创建群组失败", err)
		}
		return chatrepo.NewChatUserRepository(txStore).Create(ctx, &chatmodel.ChatUser{
			ChatId: chatEntity.Id, UserId: in.OperatorUserId, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
		})
	})
	if err != nil {
		if bizErr, ok := errs.FromError(err); ok {
			return nil, toGRPCStatus(bizErr)
		}
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "添加创建人到群组失败", err))
	}

	if len(in.UserIds) > 0 {
		l.addInitialMembers(chatEntity.Id, in.OperatorUserId, in.UserIds, now)
	}

	return &chat.Empty{}, nil
}

func (l *ChatGroupCreateLogic) addInitialMembers(chatID, creatorID uint64, userIDs []uint64, now int64) {
	existingUsers, _ := l.svcCtx.ChatUser.FindByChatID(l.ctx, chatID)
	existingUserMap := make(map[uint64]bool, len(existingUsers))
	for _, cu := range existingUsers {
		existingUserMap[cu.UserId] = true
	}

	for _, userID := range userIDs {
		if userID == creatorID || existingUserMap[userID] {
			continue
		}

		profile, err := l.svcCtx.IamCallback.GetUserProfile(l.ctx, &iamcallbackpb.GetUserProfileRequest{UserId: userID})
		if err != nil || !profile.Exists {
			logx.Errorf("查询用户失败或用户不存在: userId=%d, err=%v", userID, err)
			continue
		}

		if err := l.svcCtx.ChatUser.Create(l.ctx, &chatmodel.ChatUser{
			ChatId: chatID, UserId: userID, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			logx.Errorf("添加成员到群组失败: userId=%d, err=%v", userID, err)
			continue
		}
		existingUserMap[userID] = true
	}
}
