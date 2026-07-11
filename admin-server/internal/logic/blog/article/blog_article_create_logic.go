// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package article

import (
	"context"
	"postapocgame/admin-server/internal/dict"
	"strings"
	"time"

	"postapocgame/admin-server/internal/consts"
	"postapocgame/admin-server/internal/svc"
	"postapocgame/admin-server/internal/types"
	"postapocgame/admin-server/pkg/errs"
	jwthelper "postapocgame/admin-server/pkg/jwt"

	"postapocgame/admin-server/internal/model/blog"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlogArticleCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBlogArticleCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlogArticleCreateLogic {
	return &BlogArticleCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BlogArticleCreateLogic) BlogArticleCreate(req *types.BlogArticleCreateReq) (resp *types.Response, err error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errs.New(errs.CodeBadRequest, "文章标题不能为空")
	}

	// 从字典读取文章标题最大长度限制（默认 100 个字符）
	maxTitleLength := dict.GetIntValue(l.ctx, l.svcCtx.Repository, consts.DictCodeBlogArticleTitleMaxLength, 100)
	if err := dict.ValidateLength(title, maxTitleLength, "文章标题"); err != nil {
		return nil, err
	}

	if len(req.TagIds) == 0 {
		return nil, errs.New(errs.CodeBadRequest, "文章至少需要关联一个标签")
	}

	now := time.Now().Unix()
	article := &blog.BlogArticle{
		Title:       title,
		Content:     req.Content,
		Status:      1, // 草稿
		AuditStatus: 1, // 未提交
		Cover:       strings.TrimSpace(req.Cover),
		Summary:     strings.TrimSpace(req.Summary),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 作者信息：取当前登录后台用户
	if u, ok := jwthelper.FromContext(l.ctx); ok {
		article.AuthorId = u.UserID
		article.AuthorName = u.Username
	}

	if err = l.svcCtx.Domain.Blog.ArticleService.CreateArticle(l.ctx, article, req.TagIds); err != nil {
		return nil, err
	}

	return &types.Response{
		Code:    int(errs.CodeOK),
		Message: "创建成功",
	}, nil
}
