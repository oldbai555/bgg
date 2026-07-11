// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"context"
	"time"

	"postapocgame/admin-server/internal/repository/registry"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"

	chatmodel "postapocgame/admin-server/internal/model/chat"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatGroupCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatGroupCreateLogic {
	return &ChatGroupCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatGroupCreateLogic) ChatGroupCreate(req *types.ChatGroupCreateReq) (resp *types.Response, err error) {
	// 获取当前用户
	user, ok := jwthelper.FromContext(l.ctx)
	if !ok {
		return nil, errs.New(errs.CodeUnauthorized, "未登录或登录已过期")
	}

	// 验证群组名称
	if req.Name == "" {
		return nil, errs.New(errs.CodeBadRequest, "群组名称不能为空")
	}

	now := time.Now().Unix()

	// 1. 创建群组（type=2）+ 2. 将创建人加入群组，两条写包一个事务：
	// 避免群组建了但创建人加入失败时留下一个没有任何成员的孤儿群组。
	chatEntity := &chatmodel.Chat{
		Name:        req.Name,
		Type:        2, // 群组类型
		Avatar:      req.Avatar,
		Description: req.Description,
		CreatedBy:   user.UserID,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   0,
	}

	err = registry.Transact(l.ctx, l.svcCtx.Repository, func(ctx context.Context, txDomain *registry.Domain) error {
		if err := txDomain.Chat.Chat.Create(ctx, chatEntity); err != nil {
			return errs.Wrap(errs.CodeInternalError, "创建群组失败", err)
		}
		return txDomain.Chat.ChatUser.Create(ctx, &chatmodel.ChatUser{
			ChatId:    chatEntity.Id,
			UserId:    user.UserID,
			JoinedAt:  now,
			CreatedAt: now,
			UpdatedAt: now,
		})
	})
	if err != nil {
		if bizErr, ok := errs.FromError(err); ok {
			return nil, bizErr
		}
		return nil, errs.Wrap(errs.CodeInternalError, "添加创建人到群组失败", err)
	}

	// 3. 添加初始成员（如果提供）
	if len(req.UserIds) > 0 {
		// 验证用户是否存在且未删除
		for _, userId := range req.UserIds {
			if userId == user.UserID {
				continue // 跳过创建人
			}

			u, err := l.svcCtx.Domain.IAM.User.FindByID(l.ctx, userId)
			if err != nil {
				logx.Errorf("查询用户失败: userId=%d, err=%v", userId, err)
				continue
			}
			if u.DeletedAt != 0 {
				logx.Errorf("用户已删除: userId=%d", userId)
				continue
			}

			// 检查是否已在群组中
			existingUsers, _ := l.svcCtx.Domain.Chat.ChatUser.FindByChatID(l.ctx, chatEntity.Id)
			alreadyInGroup := false
			for _, cu := range existingUsers {
				if cu.UserId == userId {
					alreadyInGroup = true
					break
				}
			}

			if !alreadyInGroup {
				chatUser := &chatmodel.ChatUser{
					ChatId:    chatEntity.Id,
					UserId:    userId,
					JoinedAt:  now,
					CreatedAt: now,
					UpdatedAt: now,
				}
				err = l.svcCtx.Domain.Chat.ChatUser.Create(l.ctx, chatUser)
				if err != nil {
					logx.Errorf("添加成员到群组失败: userId=%d, err=%v", userId, err)
					// 继续添加其他成员，不中断流程
				}
			}
		}
	}

	return &types.Response{
		Code:    0,
		Message: "创建群组成功",
	}, nil
}
