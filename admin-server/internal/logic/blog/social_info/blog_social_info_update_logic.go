// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package social_info

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

type BlogSocialInfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogSocialInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogSocialInfoUpdateLogic {
	return &BlogSocialInfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogSocialInfoUpdateLogic) BlogSocialInfoUpdate(req *types.BlogSocialInfoUpdateReq) (resp *types.Response, err error) {
	if req.Id == 0 {
		return nil, errs.New(errs.CodeBadRequest, "社交信息ID不能为空")
	}

	// 查询现有记录
	info, err := l.svcCtx.Domain.Blog.SocialInfo.FindByID(l.ctx, req.Id)
	if err != nil {
		return nil, errs.Wrap(errs.CodeBadDB, "查询社交信息失败", err)
	}
	if info == nil {
		return nil, errs.New(errs.CodeNotFound, "社交信息不存在")
	}

	// 更新字段（仅更新提供的字段）
	if req.Name != "" {
		name := strings.TrimSpace(req.Name)
		nameMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogSocialInfoNameMaxLength, 15)
		if err := dict.ValidateLength(name, nameMaxLen, "社交平台名称"); err != nil {
			return nil, err
		}
		info.Name = name
	}

	if req.Url != "" {
		url := strings.TrimSpace(req.Url)
		urlMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogSocialInfoUrlMaxLength, 255)
		if len(url) > int(urlMaxLen) {
			return nil, errs.New(errs.CodeBadRequest, fmt.Sprintf("目标链接长度不能超过%d个字符", urlMaxLen))
		}
		info.Url = url
	}

	if req.Remark != "" {
		remark := strings.TrimSpace(req.Remark)
		remarkMaxLen := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogSocialInfoRemarkMaxLength, 127)
		if err := dict.ValidateLength(remark, remarkMaxLen, "备注"); err != nil {
			return nil, err
		}
		info.Remark = remark
	}

	if req.Status > 0 {
		if req.Status != 1 && req.Status != 2 {
			return nil, errs.New(errs.CodeBadRequest, "状态值无效，必须为1（启用）或2（禁用）")
		}
		info.Status = req.Status
	}

	if req.OrderNum > 0 {
		info.OrderNum = req.OrderNum
	}

	if err := l.svcCtx.Domain.Blog.SocialInfo.Update(l.ctx, info); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    0,
		Message: "更新成功",
	}, nil
}
