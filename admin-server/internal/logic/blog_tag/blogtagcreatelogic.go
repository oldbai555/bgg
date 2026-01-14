// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_tag

import (
	"context"
	"postapocgame/admin-server/internal/dict"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/model"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagCreateLogic {
	return &BlogTagCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogTagCreateLogic) BlogTagCreate(req *types.BlogTagCreateReq) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errs.New(errs.CodeBadRequest, "标签名称不能为空")
	}

	// 从字典读取标签名称最大长度限制（默认 10 个字符）
	maxLength := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogTagNameMaxLength, 10)
	if err := dict.ValidateLength(name, maxLength, "标签名称"); err != nil {
		return err
	}

	tag := &model.BlogTag{
		Name:   name,
		Status: req.Status,
		Remark: strings.TrimSpace(req.Remark),
	}

	if tag.Status == 0 {
		// 默认启用
		tag.Status = 1
	}

	if err := l.svcCtx.BlogTagRepository.Create(l.ctx, tag); err != nil {
		return err
	}

	return nil
}
