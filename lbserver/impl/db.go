package impl

import (
	"github.com/oldbai555/bgg/client/lbaccount"
	"github.com/oldbai555/bgg/client/lbblog"
	"github.com/oldbai555/bgg/client/lbcustomer"
	"github.com/oldbai555/bgg/client/lbim"
	"github.com/oldbai555/bgg/client/lbuser"
)

var (
	UserOrm     *OrmCondBuilder
	ArticleOrm  *OrmCondBuilder
	CategoryOrm *OrmCondBuilder
	CommentOrm  *OrmCondBuilder
	MessageOrm  *OrmCondBuilder
	CustomerOrm *OrmCondBuilder
	AccountOrm  *OrmCondBuilder
)

func InitDbOrm() {
	UserOrm = NewOrmCondBuilder(
		&lbuser.ModelUser{},
		lbuser.ErrUserNotFound,
	)
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

	MessageOrm = NewOrmCondBuilder(
		&lbim.ModelMessage{},
		lbim.ErrNotFoundMessage,
	)

	CustomerOrm = NewOrmCondBuilder(
		&lbcustomer.ModelCustomer{},
		lbcustomer.ErrNotFoundCustomer,
	)

	AccountOrm = NewOrmCondBuilder(
		&lbaccount.ModelAccount{},
		lbaccount.ErrNotFoundAccount,
	)

}
