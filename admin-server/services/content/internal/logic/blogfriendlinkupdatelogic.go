package logic

import (
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogFriendLinkUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogFriendLinkUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkUpdateLogic {
	return &BlogFriendLinkUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogFriendLinkUpdate 迁移自 internal/logic/blog/friend_link/blog_friend_link_update_logic.go。
func (l *BlogFriendLinkUpdateLogic) BlogFriendLinkUpdate(in *content.BlogFriendLinkUpdateRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "链接ID不能为空"))
	}

	link, err := l.svcCtx.FriendLink.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询友情链接失败", err))
	}
	if link == nil {
		return nil, toGRPCStatus(errs.New(errs.CodeNotFound, "友情链接不存在"))
	}

	limits := l.svcCtx.Config.Limits

	if in.Name != "" {
		name := strings.TrimSpace(in.Name)
		if err := validateLength(name, limits.BlogFriendLinkNameMaxLength, "链接名称"); err != nil {
			return nil, toGRPCStatus(err)
		}
		link.Name = name
	}
	if in.Url != "" {
		url := strings.TrimSpace(in.Url)
		if int64(len([]rune(url))) > limits.BlogFriendLinkUrlMaxLength {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, fmt.Sprintf("目标链接长度不能超过%d个字符", limits.BlogFriendLinkUrlMaxLength)))
		}
		link.Url = url
	}
	if in.Remark != "" {
		remark := strings.TrimSpace(in.Remark)
		if err := validateLength(remark, limits.BlogFriendLinkRemarkMaxLength, "备注"); err != nil {
			return nil, toGRPCStatus(err)
		}
		link.Remark = remark
	}
	if in.Status > 0 {
		if in.Status != 1 && in.Status != 2 {
			return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "状态值无效，必须为1（启用）或2（禁用）"))
		}
		link.Status = in.Status
	}
	if in.OrderNum > 0 {
		link.OrderNum = in.OrderNum
	}

	if err := l.svcCtx.FriendLink.Update(l.ctx, link); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
