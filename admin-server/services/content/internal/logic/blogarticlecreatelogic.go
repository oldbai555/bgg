package logic

import (
	"context"
	"strings"
	"time"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/content"
	contentconsts "postapocgame/admin-server/services/content/internal/consts"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
	"postapocgame/admin-server/services/content/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlogArticleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleCreateLogic {
	return &BlogArticleCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BlogArticleCreate 迁移自 internal/logic/blog/article/blog_article_create_logic.go。
// 标题长度上限原来读字典 blog_article_title_max_length（物理属于 iam 域），改成
// svcCtx.Config.Limits.BlogArticleTitleMaxLength 静态配置。
func (l *BlogArticleCreateLogic) BlogArticleCreate(in *content.BlogArticleCreateRequest) (*content.Empty, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章标题不能为空"))
	}
	if err := validateLength(title, l.svcCtx.Config.Limits.BlogArticleTitleMaxLength, "文章标题"); err != nil {
		return nil, toGRPCStatus(err)
	}
	if len(in.TagIds) == 0 {
		return nil, toGRPCStatus(errs.New(errs.CodeBadRequest, "文章至少需要关联一个标签"))
	}

	now := time.Now().Unix()
	article := &blogmodel.BlogArticle{
		Title:       title,
		Content:     in.Content,
		Status:      contentconsts.BlogArticleStatusDraft,
		AuditStatus: contentconsts.BlogArticleAuditStatusNotSubmitted,
		Cover:       strings.TrimSpace(in.Cover),
		Summary:     strings.TrimSpace(in.Summary),
		AuthorId:    in.OperatorUserId,
		AuthorName:  in.OperatorUsername,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := l.svcCtx.ArticleService.CreateArticle(l.ctx, article, in.TagIds); err != nil {
		return nil, toGRPCStatus(err)
	}

	return &content.Empty{}, nil
}
