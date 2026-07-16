package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogTagUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagUpdateLogic {
	return &BlogTagUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogTagUpdate 迁移自 internal/logic/blog/tag/blog_tag_update_logic.go。
func (l *BlogTagUpdateLogic) BlogTagUpdate(in *content.BlogTagUpdateRequest) (*content.Empty, error) {
	if in.Id == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "标签ID不能为空"))
	}

	tag, err := l.svcCtx.Tag.FindByID(l.ctx, in.Id)
	if err != nil {
		return nil, toGRPCStatus(errs.Wrap(errs.CodeBadDB, "查询标签失败", err))
	}

	if in.Name != "" {
		name := strings.TrimSpace(in.Name)
		if err := validateLength(name, l.svcCtx.Config.Limits.BlogTagNameMaxLength, "标签名称"); err != nil {
			return nil, toGRPCStatus(err)
		}
		tag.Name = name
	}
	if in.Status != 0 {
		tag.Status = in.Status
	}
	if in.Remark != "" {
		tag.Remark = strings.TrimSpace(in.Remark)
	}

	if err := l.svcCtx.Tag.Update(l.ctx, tag); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
