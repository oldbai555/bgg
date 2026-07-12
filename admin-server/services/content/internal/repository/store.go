// Package repository 从 internal/repository/blog/、internal/repository/video/ 原样搬迁而来。
// 唯一的结构性改动是把各 repository 原来共享的 *repository.Repository（单体聚合了全部
// 业务域 Model 的大句柄）换成这里的 *Store——content-rpc 从第一天起只有 blog 六张表
// （blog_article/blog_article_tag/blog_article_audit/blog_friend_link/blog_social_info/
// blog_tag）+ video 一张表，不该也不能继续持有指向其它域的句柄。和 services/sdk、
// services/chat 的 Store 同一个模式。
package repository

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	blogmodel "postapocgame/admin-server/services/content/internal/model/blog"
	videomodel "postapocgame/admin-server/services/content/internal/model/video"
)

// Store 聚合 content-rpc 自己需要的全部 Model，供 blog/video 各 repository 共用；
// BlogArticleService（services/content/internal/domain/content）需要事务时直接调
// Store.Transact。
type Store struct {
	DB                    sqlx.SqlConn
	BlogArticleModel      blogmodel.BlogArticleModel
	BlogArticleTagModel   blogmodel.BlogArticleTagModel
	BlogArticleAuditModel blogmodel.BlogArticleAuditModel
	BlogFriendLinkModel   blogmodel.BlogFriendLinkModel
	BlogSocialInfoModel   blogmodel.BlogSocialInfoModel
	BlogTagModel          blogmodel.BlogTagModel
	VideoModel            videomodel.VideoModel
}

func NewStore(conn sqlx.SqlConn, cacheConf cache.CacheConf) *Store {
	return &Store{
		DB:                    conn,
		BlogArticleModel:      blogmodel.NewBlogArticleModel(conn, cacheConf),
		BlogArticleTagModel:   blogmodel.NewBlogArticleTagModel(conn, cacheConf),
		BlogArticleAuditModel: blogmodel.NewBlogArticleAuditModel(conn, cacheConf),
		BlogFriendLinkModel:   blogmodel.NewBlogFriendLinkModel(conn, cacheConf),
		BlogSocialInfoModel:   blogmodel.NewBlogSocialInfoModel(conn, cacheConf),
		BlogTagModel:          blogmodel.NewBlogTagModel(conn, cacheConf),
		VideoModel:            videomodel.NewVideoModel(conn, cacheConf),
	}
}

// Transact 在单个 MySQL 事务内执行 fn，用法和 internal/repository/repository.go 的
// Repository.Transact 完全同构（content-rpc 自己的小号版本）。
func (s *Store) Transact(ctx context.Context, fn func(ctx context.Context, txStore *Store) error) error {
	return s.DB.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, s.withSession(session))
	})
}

func (s *Store) withSession(session sqlx.Session) *Store {
	return &Store{
		DB:                    sqlx.NewSqlConnFromSession(session),
		BlogArticleModel:      s.BlogArticleModel.WithSession(session),
		BlogArticleTagModel:   s.BlogArticleTagModel.WithSession(session),
		BlogArticleAuditModel: s.BlogArticleAuditModel.WithSession(session),
		BlogFriendLinkModel:   s.BlogFriendLinkModel.WithSession(session),
		BlogSocialInfoModel:   s.BlogSocialInfoModel.WithSession(session),
		BlogTagModel:          s.BlogTagModel.WithSession(session),
		VideoModel:            s.VideoModel.WithSession(session),
	}
}
