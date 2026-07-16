package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	iamcallbackpb "postapocgame/admin-server/pkg/iamcallback/pb"
	"postapocgame/admin-server/services/chat/chat"
	chatconsts "postapocgame/admin-server/services/chat/internal/consts"
	chatmodel "postapocgame/admin-server/services/chat/internal/model/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupDetailLogic {
	return &ChatGroupDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupDetail 迁移自 internal/logic/chat/group/chat_group_detail_logic.go。
func (l *ChatGroupDetailLogic) ChatGroupDetail(in *chat.ChatGroupDetailRequest) (*chat.ChatGroupDetailResponse, error) {
	c, err := l.svcCtx.Chat.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeNotFound, "群组不存在", err))
	}
	if c.Type != chatconsts.ChatTypeGroup {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "该聊天不是群组"))
	}
	if c.DeletedAt != 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "群组已删除"))
	}

	chatUsers, err := l.svcCtx.ChatUser.FindByChatID(l.ctx, c.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询群组成员失败", err))
	}

	members := resolveGroupMembers(l.ctx, l.svcCtx, chatUsers)

	return &chat.ChatGroupDetailResponse{
		Id:          c.Id,
		Name:        c.Name,
		Avatar:      c.Avatar,
		Description: c.Description,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		MemberCount: int64(len(members)),
		Members:     members,
	}, nil
}

// resolveGroupMembers 迁移自 internal/logic/chat/group/{chat_group_detail_logic.go,
// chat_group_member_list_logic.go} 里重复出现的"逐个用户回调 IamCallback.GetUserProfile
// 拿用户名/昵称/头像/部门名/角色名"逻辑，ChatGroupDetail 和 ChatGroupMemberList 共用同一份
// 实现（和拆分前两个 logic 文件各自内联一份完全相同代码相比，这里去掉了重复）。
func resolveGroupMembers(ctx context.Context, svcCtx *svc.ServiceContext, chatUsers []chatmodel.ChatUser) []*chat.ChatGroupMemberItem {
	members := make([]*chat.ChatGroupMemberItem, 0, len(chatUsers))
	for _, cu := range chatUsers {
		profile, err := svcCtx.IamCallback.GetUserProfile(ctx, &iamcallbackpb.GetUserProfileRequest{UserId: cu.UserId})
		if err != nil || !profile.Exists {
			continue // 跳过已删除/查询失败的用户
		}

		members = append(members, &chat.ChatGroupMemberItem{
			UserId:         cu.UserId,
			Username:       profile.Username,
			Nickname:       profile.Nickname,
			Avatar:         profile.Avatar,
			DepartmentName: profile.DepartmentName,
			RoleNames:      profile.RoleNames,
			JoinedAt:       cu.JoinedAt,
		})
	}
	return members
}
