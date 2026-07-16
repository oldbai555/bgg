package content

import (
	"context"

	"postapocgame/admin-server/pkg/errs"
	"postapocgame/admin-server/services/content/internal/consts"
	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
	"postapocgame/admin-server/services/content/internal/repository"
	blogrepo "postapocgame/admin-server/services/content/internal/repository/blog"
)

// BlogArticleService 承载文章创建/更新（跨 blog_article + blog_article_tag 两表写）
// 和审核/下架（跨 blog_article_audit + blog_article 两表写）。
type BlogArticleService struct {
	store *repository.Store
}

func NewBlogArticleService(store *repository.Store) *BlogArticleService {
	return &BlogArticleService{store: store}
}

func (s *BlogArticleService) CreateArticle(ctx context.Context, article *blogmodel.BlogArticle, tagIDs []uint64) error {
	return s.store.Transact(ctx, func(ctx context.Context, txStore *repository.Store) error {
		return blogrepo.NewBlogArticleRepository(txStore).CreateWithTags(ctx, article, tagIDs)
	})
}

func (s *BlogArticleService) UpdateArticle(ctx context.Context, article *blogmodel.BlogArticle, tagIDs []uint64) error {
	return s.store.Transact(ctx, func(ctx context.Context, txStore *repository.Store) error {
		return blogrepo.NewBlogArticleRepository(txStore).UpdateWithTags(ctx, article, tagIDs)
	})
}

// AuditArticle 对应 blog_article_audit_logic.go 现有逻辑：写审核记录 + 更新文章审核状态，
// 原来是两次独立写，这里包进同一个事务。
func (s *BlogArticleService) AuditArticle(ctx context.Context, articleID uint64, result int64, remark string, auditorID uint64, auditorName string) (*blogmodel.BlogArticle, error) {
	var article *blogmodel.BlogArticle
	err := s.store.Transact(ctx, func(ctx context.Context, txStore *repository.Store) error {
		articleRepo := blogrepo.NewBlogArticleRepository(txStore)
		a, err := articleRepo.FindByID(ctx, articleID)
		if err != nil {
			return errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
		}
		if a == nil || a.DeletedAt != 0 {
			return errs.New(errs.CodeNotFound, "文章不存在")
		}
		if a.AuditStatus != consts.BlogArticleAuditStatusPending {
			return errs.New(errs.CodeForbidden, "当前状态不允许审核")
		}

		if err := blogrepo.NewBlogArticleAuditRepository(txStore).Create(ctx, &blogmodel.BlogArticleAudit{
			ArticleId: a.Id, AuditStatus: result, AuditRemark: remark,
			AuditorId: auditorID, AuditorName: auditorName,
		}); err != nil {
			return err
		}

		a.AuditStatus = result
		if result == consts.BlogArticleAuditStatusPassed {
			a.Status = consts.BlogArticleStatusAuditPassed
		}
		if err := articleRepo.Update(ctx, a); err != nil {
			return errs.Wrap(errs.CodeBadDB, "更新文章审核状态失败", err)
		}
		article = a
		return nil
	})
	return article, err
}

// SetArticleTop 对应 blog_article_top_logic.go 现有逻辑：若已达置顶数量上限，先取消最早置顶的
// 文章再置顶目标文章，两次 UpdateTopStatus 写原来是独立的、无事务保护，这里包进同一个事务，
// 避免中途失败导致临时出现 0 篇置顶文章。
func (s *BlogArticleService) SetArticleTop(ctx context.Context, articleID uint64, maxCount int64) error {
	return s.store.Transact(ctx, func(ctx context.Context, txStore *repository.Store) error {
		articleRepo := blogrepo.NewBlogArticleRepository(txStore)

		article, err := articleRepo.FindByID(ctx, articleID)
		if err != nil {
			return errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
		}
		if article == nil || article.DeletedAt != 0 {
			return errs.New(errs.CodeNotFound, "文章不存在")
		}
		if article.IsTop == 1 {
			return nil
		}

		currentTopCount, err := articleRepo.FindTopCount(ctx)
		if err != nil {
			return errs.Wrap(errs.CodeBadDB, "查询置顶文章数量失败", err)
		}
		if currentTopCount >= maxCount {
			oldestTopArticle, err := articleRepo.FindOldestTopArticle(ctx)
			if err != nil {
				return errs.Wrap(errs.CodeBadDB, "查询最早置顶文章失败", err)
			}
			if oldestTopArticle != nil {
				if err := articleRepo.UpdateTopStatus(ctx, oldestTopArticle.Id, 0); err != nil {
					return errs.Wrap(errs.CodeBadDB, "取消最早置顶文章失败", err)
				}
			}
		}

		if err := articleRepo.UpdateTopStatus(ctx, articleID, 1); err != nil {
			return errs.Wrap(errs.CodeBadDB, "设置文章置顶失败", err)
		}
		return nil
	})
}

// UnpublishArticle 对应 blog_article_audit_unpublish_logic.go 现有逻辑：更新文章下架状态 +
// 写一条审核记录（便于追踪谁下架），与 AuditArticle 同一模式，包进同一个事务。
func (s *BlogArticleService) UnpublishArticle(ctx context.Context, articleID uint64, remark string, operatorID uint64, operatorName string) (*blogmodel.BlogArticle, error) {
	var article *blogmodel.BlogArticle
	err := s.store.Transact(ctx, func(ctx context.Context, txStore *repository.Store) error {
		articleRepo := blogrepo.NewBlogArticleRepository(txStore)
		a, err := articleRepo.FindByID(ctx, articleID)
		if err != nil {
			return errs.Wrap(errs.CodeBadDB, "查询文章失败", err)
		}
		if a == nil || a.DeletedAt != 0 {
			return errs.New(errs.CodeNotFound, "文章不存在")
		}
		if a.Status != consts.BlogArticleStatusPublished {
			return errs.New(errs.CodeForbidden, "仅已上架文章可下架")
		}

		a.Status = consts.BlogArticleStatusUnpublished
		if err := articleRepo.Update(ctx, a); err != nil {
			return errs.Wrap(errs.CodeBadDB, "下架失败", err)
		}

		if err := blogrepo.NewBlogArticleAuditRepository(txStore).Create(ctx, &blogmodel.BlogArticleAudit{
			ArticleId:   a.Id,
			AuditStatus: consts.BlogArticleAuditStatusRejected, // 复用字段，语义为"审核操作类记录"
			AuditRemark: remark,
			AuditorId:   operatorID,
			AuditorName: operatorName,
		}); err != nil {
			return errs.Wrap(errs.CodeBadDB, "记录下架操作失败", err)
		}
		article = a
		return nil
	})
	return article, err
}
