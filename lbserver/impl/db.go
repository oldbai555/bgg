package impl

import (
	"github.com/oldbai555/bgg/lbblog"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	UserOrm     *OrmCondBuilder
	ArticleOrm  *OrmCondBuilder
	CategoryOrm *OrmCondBuilder
	CommentOrm  *OrmCondBuilder
)

func InitDbOrm() {
	UserOrm = NewOrmCondBuilder(
		&lbuser.ModelUser{},
		lberr.NewErr(int32(lbuser.ErrCode_ErrUserNotFound), "ErrUserNotFound"),
	)
	ArticleOrm = NewOrmCondBuilder(
		&lbblog.ModelArticle{},
		lberr.NewErr(int32(lbblog.ErrCode_ErrArticleNotFound), "ErrArticleNotFound"),
	)

	CategoryOrm = NewOrmCondBuilder(
		&lbblog.ModelCategory{},
		lberr.NewErr(int32(lbblog.ErrCode_ErrCategoryNotFound), "ErrCategoryNotFound"),
	)

	CommentOrm = NewOrmCondBuilder(
		&lbblog.ModelComment{},
		lberr.NewErr(int32(lbblog.ErrCode_ErrCommentNotFound), "ErrCommentNotFound"),
	)
}
