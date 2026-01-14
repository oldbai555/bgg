// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package blog_tag

import (
	"context"
	"postapocgame/admin-server/internal/dict"
	"strings"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogTagUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagUpdateLogic {
	return &BlogTagUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogTagUpdateLogic) BlogTagUpdate(req *types.BlogTagUpdateReq) error {
	if req.Id == 0 {
		return errs.New(errs.CodeBadRequest, "标签ID不能为空")
	}

	tagRepo := l.svcCtx.BlogTagRepository

	tag, err := tagRepo.FindByID(l.ctx, req.Id)
	if err != nil {
		return errs.Wrap(errs.CodeBadDB, "查询标签失败", err)
	}

	if req.Name != "" {
		name := strings.TrimSpace(req.Name)
		// 从字典读取标签名称最大长度限制（默认 10 个字符）
		maxLength := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogTagNameMaxLength, 10)
		if err := dict.ValidateLength(name, maxLength, "标签名称"); err != nil {
			return err
		}
		tag.Name = name
	}
	if req.Status != 0 {
		tag.Status = req.Status
	}
	if req.Remark != "" {
		tag.Remark = strings.TrimSpace(req.Remark)
	}

	if err = tagRepo.Update(l.ctx, tag); err != nil {
		return err
	}

	return nil
}
