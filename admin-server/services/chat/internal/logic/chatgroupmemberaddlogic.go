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
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupMemberAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupMemberAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupMemberAddLogic {
	return &ChatGroupMemberAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupMemberAdd 迁移自 internal/logic/chat/group/chat_group_member_add_logic.go。
func (l *ChatGroupMemberAddLogic) ChatGroupMemberAdd(in *chat.ChatGroupMemberAddRequest) (*chat.Empty, error) {
	if len(in.UserIds) == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "用户ID列表不能为空"))
	}

	c, err := l.svcCtx.Chat.FindByID(l.ctx, in.ChatId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "群组不存在", err))
	}
	if c.Type != chatconsts.ChatTypeGroup {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该聊天不是群组"))
	}
	if c.DeletedAt != 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "群组已删除"))
	}

	existingUsers, err := l.svcCtx.ChatUser.FindByChatID(l.ctx, in.ChatId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询群组成员失败", err))
	}
	existingUserMap := make(map[uint64]bool, len(existingUsers))
	for _, cu := range existingUsers {
		existingUserMap[cu.UserId] = true
	}

	now := time.Now().Unix()
	addedCount := 0
	for _, userID := range in.UserIds {
		if existingUserMap[userID] {
			continue
		}

		profile, err := l.svcCtx.IamCallback.GetUserProfile(l.ctx, &iamcallbackpb.GetUserProfileRequest{UserId: userID})
		if err != nil || !profile.Exists {
			logx.Errorf("查询用户失败或用户不存在: userId=%d, err=%v", userID, err)
			continue
		}

		if err := l.svcCtx.ChatUser.Create(l.ctx, &chatmodel.ChatUser{
			ChatId: in.ChatId, UserId: userID, JoinedAt: now, CreatedAt: now, UpdatedAt: now,
		}); err != nil {
			logx.Errorf("添加成员到群组失败: userId=%d, err=%v", userID, err)
			continue
		}
		addedCount++
		existingUserMap[userID] = true
	}

	if addedCount == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "没有可添加的成员（可能已全部在群组中或用户不存在）"))
	}

	return &chat.Empty{}, nil
}
