package impl

import (
	"github.com/oldbai555/bgg/client/lbblog"
)

var (
	ArticleOrm  *OrmCondBuilder
	CategoryOrm *OrmCondBuilder
	CommentOrm  *OrmCondBuilder
)

func InitDbOrm() {

	ArticleOrm = NewOrmCondBuilder(
		&lbblog.ModelArticle{},
		lbblog.ErrArticleNotFound,
	)

	CategoryOrm = NewOrmCondBuilder(
		&lbblog.ModelCategory{},
		lbblog.ErrCategoryNotFound,
	)

	CommentOrm = NewOrmCondBuilder(
		&lbblog.ModelComment{},
		lbblog.ErrCommentNotFound,
	)

}
