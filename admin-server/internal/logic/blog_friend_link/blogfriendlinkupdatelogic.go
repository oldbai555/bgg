// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_friend_link

import (
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/dict"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogFriendLinkUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogFriendLinkUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkUpdateLogic {
	return &BlogFriendLinkUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogFriendLinkUpdateLogic) BlogFriendLinkUpdate(req *types.BlogFriendLinkUpdateReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "链接ID不能为空")
	}

	// 查询现有记录
	link, err := l.svcCtx.BlogFriendLinkRepository.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询友情链接失败", err)
	}
	if link == nil {
		return nil, errs.New(errs.CodeNotFound, "友情链接不存在")
	}

	// 更新字段（仅更新提供的字段）
	if req.Name != "" {
		name := strings.TrimSpace(req.Name)
		nameMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogFriendLinkNameMaxLength, 15)
		if err := dict.ValidateLength(name, nameMaxLen, "链接名称"); err != nil {
			return nil, err
		}
		link.Name = name
	}

	if req.Url != "" {
		url := strings.TrimSpace(req.Url)
		urlMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogFriendLinkUrlMaxLength, 255)
		if len(url) > int(urlMaxLen) {
			return nil, errs.New(errs.CodeBadRequest, fmt.Sprintf("目标链接长度不能超过%d个字符", urlMaxLen))
		}
		link.Url = url
	}

	if req.Remark != "" {
		remark := strings.TrimSpace(req.Remark)
		remarkMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogFriendLinkRemarkMaxLength, 127)
		if err := dict.ValidateLength(remark, remarkMaxLen, "备注"); err != nil {
			return nil, err
		}
		link.Remark = remark
	}

	if req.Status > 0 {
		if req.Status != 1 && req.Status != 2 {
			return nil, errs.New(errs.CodeBadRequest, "状态值无效，必须为1（启用）或2（禁用）")
		}
		link.Status = req.Status
	}

	if req.OrderNum > 0 {
		link.OrderNum = req.OrderNum
	}

	if err := l.svcCtx.BlogFriendLinkRepository.Update(l.ctx, link); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "更新成功",
	}, nil
}
