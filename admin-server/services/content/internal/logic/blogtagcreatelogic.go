package logic

import (
	"context"
	"strings"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogTagCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogTagCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogTagCreateLogic {
	return &BlogTagCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogTagCreate 迁移自 internal/logic/blog/tag/blog_tag_create_logic.go。
func (l *BlogTagCreateLogic) BlogTagCreate(in *content.BlogTagCreateRequest) (*content.Empty, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "标签名称不能为空"))
	}
	if err := validateLength(name, l.svcCtx.Config.Limits.BlogTagNameMaxLength, "标签名称"); err != nil {
		return nil, toGRPCStatus(err)
	}

	tag := &blogmodel.BlogTag{
		Name:   name,
		Status: in.Status,
		Remark: strings.TrimSpace(in.Remark),
	}
	if tag.Status == 0 {
		tag.Status = 1 // 默认启用
	}

	if err := l.svcCtx.Tag.Create(l.ctx, tag); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
