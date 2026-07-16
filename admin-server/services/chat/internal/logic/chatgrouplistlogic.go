package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/chat/chat"
	"postapocgame/admin-server/services/chat/internal/svc"
)

type ChatGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupListLogic {
	return &ChatGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChatGroupList 迁移自 internal/logic/chat/group/chat_group_list_logic.go。
func (l *ChatGroupListLogic) ChatGroupList(in *chat.ChatGroupListRequest) (*chat.ChatGroupListResponse, error) {
	page, pageSize := normalizePage(in.Page, in.PageSize, 10, 100)

	groups, total, err := l.svcCtx.Chat.FindGroups(l.ctx, page, pageSize, in.Name)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeInternalError, "查询群组列表失败", err))
	}

	items := make([]*chat.ChatGroupItem, 0, len(groups))
	for _, group := range groups {
		memberCount, err := l.svcCtx.Chat.CountMembersByChatID(l.ctx, group.Id)
		if err != nil {
			logx.Errorf("统计群组 %d 成员数量失败: %v", group.Id, err)
			memberCount = 0
		}

		items = append(items, &chat.ChatGroupItem{
			Id:          group.Id,
			Name:        group.Name,
			Avatar:      group.Avatar,
			Description: group.Description,
			CreatedBy:   group.CreatedBy,
			CreatedAt:   group.CreatedAt,
			MemberCount: memberCount,
		})
	}

	return &chat.ChatGroupListResponse{Total: total, List: items}, nil
}
