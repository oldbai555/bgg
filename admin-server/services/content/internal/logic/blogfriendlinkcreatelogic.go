package logic

import (
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogFriendLinkCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogFriendLinkCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkCreateLogic {
	return &BlogFriendLinkCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogFriendLinkCreate 迁移自 internal/logic/blog/friend_link/blog_friend_link_create_logic.go。
// 长度上限原来读字典（物理属于 iam 域），改成 svcCtx.Config.Limits 静态配置。
func (l *BlogFriendLinkCreateLogic) BlogFriendLinkCreate(in *content.BlogFriendLinkCreateRequest) (*content.Empty, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "链接名称不能为空"))
	}
	url := strings.TrimSpace(in.Url)
	if url == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "目标链接不能为空"))
	}

	limits := l.svcCtx.Config.Limits
	if err := validateLength(name, limits.BlogFriendLinkNameMaxLength, "链接名称"); err != nil {
		return nil, toGRPCStatus(err)
	}
	if int64(len([]rune(url))) > limits.BlogFriendLinkUrlMaxLength {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, fmt.Sprintf("目标链接长度不能超过%d个字符", limits.BlogFriendLinkUrlMaxLength)))
	}

	remark := strings.TrimSpace(in.Remark)
	if remark != "" {
		if err := validateLength(remark, limits.BlogFriendLinkRemarkMaxLength, "备注"); err != nil {
			return nil, toGRPCStatus(err)
		}
	}

	if in.Status != 1 && in.Status != 2 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "状态值无效，必须为1（启用）或2（禁用）"))
	}

	link := &blogmodel.BlogFriendLink{
		Name:     name,
		Url:      url,
		Remark:   remark,
		Status:   in.Status,
		OrderNum: in.OrderNum,
	}
	if err := l.svcCtx.FriendLink.Create(l.ctx, link); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
