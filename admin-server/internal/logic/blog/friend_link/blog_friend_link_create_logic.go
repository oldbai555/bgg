// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package friend_link

import (
	blogrepo "postapocgame/admin-server/internal/repository/blog"
	"context"
	"fmt"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/dict"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
	"postapocgame/admin-server/internal/model/blog"
)

type BlogFriendLinkCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogFriendLinkCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogFriendLinkCreateLogic {
	return &BlogFriendLinkCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogFriendLinkCreateLogic) BlogFriendLinkCreate(req *types.BlogFriendLinkCreateReq) (resp *types.Response, err error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errs.New(errs.CodeBadRequest, "链接名称不能为空")
	}

	url := strings.TrimSpace(req.Url)
	if url == "" {
		return nil, errs.New(errs.CodeBadRequest, "目标链接不能为空")
	}

	// 从字典读取长度限制并校验
	nameMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogFriendLinkNameMaxLength, 15)
	if err := dict.ValidateLength(name, nameMaxLen, "链接名称"); err != nil {
		return nil, err
	}

	urlMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogFriendLinkUrlMaxLength, 255)
	if len(url) > int(urlMaxLen) {
		return nil, errs.New(errs.CodeBadRequest, fmt.Sprintf("目标链接长度不能超过%d个字符", urlMaxLen))
	}

	remark := strings.TrimSpace(req.Remark)
	if remark != "" {
		remarkMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogFriendLinkRemarkMaxLength, 127)
		if err := dict.ValidateLength(remark, remarkMaxLen, "备注"); err != nil {
			return nil, err
		}
	}

	// 状态校验
	if req.Status != 1 && req.Status != 2 {
		return nil, errs.New(errs.CodeBadRequest, "状态值无效，必须为1（启用）或2（禁用）")
	}

	link := &blog.BlogFriendLink{
		Name:     name,
		Url:      url,
		Remark:   remark,
		Status:   req.Status,
		OrderNum: req.OrderNum,
	}

	if err := blogrepo.NewBlogFriendLinkRepository(l.svcCtx.Repository).Create(l.ctx, link); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "创建成功",
	}, nil
}
