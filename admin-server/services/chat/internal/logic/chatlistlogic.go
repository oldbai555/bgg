package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/chat/chat"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatListLogic {
	return &ChatListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatList 迁移自 internal/logic/chat/chat/chat_list_logic.go：私聊场景下对方的用户名/昵称/
// 部门名/角色名列表原来是本地直接读 Domain.IAM.*，现在改成回调单体内嵌的
// IamCallback.GetUserProfile（见 pkg/iamcallback 包注释）。
func (l *ChatListLogic) ChatList(in *chat.ChatListRequest) (*chat.ChatListResponse, error) {
	if in.OperatorUserId == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeUnauthorized, "未登录或登录已过期"))
	}

	chats, err := l.svcCtx.Chat.FindByUserID(l.ctx, in.OperatorUserId)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询聊天列表失败", err))
	}

	items := make([]*chat.ChatItem, 0, len(chats))
	for _, c := range chats {
		item := &chat.ChatItem{
			ChatId:      c.Id,
			Name:        c.Name,
			ChatType:    c.Type,
			Avatar:      c.Avatar,
			Description: c.Description,
		}

		if c.Type == chatconsts.ChatTypePrivate {
			chatUsers, err := l.svcCtx.ChatUser.FindByChatID(l.ctx, c.Id)
			if err == nil && len(chatUsers) == 2 {
				var otherUserID uint64
				for _, cu := range chatUsers {
					if cu.UserId != in.OperatorUserId {
						otherUserID = cu.UserId
						break
					}
				}
				if otherUserID > 0 {
					l.fillOtherUserProfile(item, otherUserID)
				}
			}
		}

		items = append(items, item)
	}

	return &chat.ChatListResponse{List: items}, nil
}

func (l *ChatListLogic) fillOtherUserProfile(item *chat.ChatItem, otherUserID uint64) {
	profile, err := l.svcCtx.IamCallback.GetUserProfile(l.ctx, &iamcallbackpb.GetUserProfileRequest{UserId: otherUserID})
	if err != nil || !profile.Exists {
		return
	}
	item.UserId = otherUserID
	item.Username = profile.Username
	item.Nickname = profile.Nickname
	if profile.Nickname != "" {
		item.Name = profile.Nickname
	} else {
		item.Name = profile.Username
	}
	item.Avatar = profile.Avatar
	item.DepartmentName = profile.DepartmentName
	item.RoleNames = profile.RoleNames
}
