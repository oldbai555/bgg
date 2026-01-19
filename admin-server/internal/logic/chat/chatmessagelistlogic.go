// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package chat

import (
	"context"
	"postapocgame/admin-server/internal/logic/logicutil"
	"postapocgame/admin-server/internal/repository"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatMessageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatMessageListLogic {
	return &ChatMessageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatMessageListLogic) ChatMessageList(req *types.ChatMessageListReq) (resp *types.ChatMessageListResp, err error) {
	if req == nil {
		return nil, errs.New(errs.CodeBadRequest, "请求参数不能为空")
	}

	chatMessageLimit := l.getChatMessageLimitFromCache()
	req.Page, req.PageSize = logicutil.NormalizePage(req.Page, req.PageSize, chatMessageLimit, chatMessageLimit)

	messageRepo := repository.NewChatMessageRepository(l.svcCtx.Repository)
	userRepo := repository.NewUserRepository(l.svcCtx.Repository)

	// 根据 chatId 查询消息
	if req.ChatId == 0 {
		return nil, errs.New(errs.CodeBadRequest, "chatId 不能为空")
	}

	list, total, err := messageRepo.FindByChatID(l.ctx, req.Page, req.PageSize, req.ChatId)
	if err != nil {
		return nil, errs.Wrap(errs.CodeInternalError, "查询聊天消息列表失败", err)
	}

	items := make([]types.ChatMessageItem, 0, len(list))
	for _, msg := range list {
		// 查询发送用户信息
		fromUser, _ := userRepo.FindByID(l.ctx, msg.FromUserId)
		fromUserName := ""
		if fromUser != nil {
			fromUserName = fromUser.Username
		}

		items = append(items, types.ChatMessageItem{
			Id:           msg.Id,
			ChatId:       msg.ChatId,
			FromUserId:   msg.FromUserId,
			FromUserName: fromUserName,
			Content:      msg.Content,
			MessageType:  msg.MessageType,
			Status:       msg.Status,
			CreatedAt:    msg.CreatedAt,
		})
	}

	return &types.ChatMessageListResp{
		Total: total,
		List:  items,
	}, nil
}

// getChatMessageLimitFromCache 读取聊天消息数量限制，优先使用缓存，兜底字典
func (l *ChatMessageListLogic) getChatMessageLimitFromCache() int64 {
	const (
		cacheKey        = "chat:config:message_limit"
		defaultLimit    = int64(30)
		cacheExpireSecs = 600 // 10 分钟
	)

	// 尝试从业务缓存读取
	var cached int64
	if err := l.svcCtx.Repository.BusinessCache.Get(l.ctx, cacheKey, &cached); err == nil && cached > 0 {
		return cached
	}

	limit := defaultLimit
	dictTypeRepo := repository.NewDictTypeRepository(l.svcCtx.Repository)
	dictType, err := dictTypeRepo.FindByCode(l.ctx, "chat_config")
	if err == nil && dictType != nil {
		dictItemRepo := repository.NewDictItemRepository(l.svcCtx.Repository)
		items, err := dictItemRepo.FindByTypeID(l.ctx, dictType.Id)
		if err == nil {
			for _, item := range items {
				if item.Label == "聊天窗口消息数量" && item.Value != "" {
					if v, parseErr := strconv.ParseInt(item.Value, 10, 64); parseErr == nil && v > 0 {
						limit = v
						break
					}
				}
			}
		}
	}

	// 缓存结果，避免频繁查字典
	_ = l.svcCtx.Repository.BusinessCache.Set(l.ctx, cacheKey, limit, cacheExpireSecs+int(time.Now().Unix()%60))

	return limit
}
