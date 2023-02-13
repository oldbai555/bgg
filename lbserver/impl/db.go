package impl

import (
	"github.com/oldbai555/bgg/lbblog"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/bgg/lbwebsocket"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	UserOrm     *OrmCondBuilder
	ArticleOrm  *OrmCondBuilder
	CategoryOrm *OrmCondBuilder
	CommentOrm  *OrmCondBuilder
	ChatOrm     *OrmCondBuilder
	MessageOrm  *OrmCondBuilder
	VisitorOrm  *OrmCondBuilder
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

	ChatOrm = NewOrmCondBuilder(
		&lbwebsocket.ModelChat{},
		lberr.NewErr(int32(lbwebsocket.ErrCode_ErrChatNotFound), "ErrChatNotFound"),
	)

	MessageOrm = NewOrmCondBuilder(
		&lbwebsocket.ModelMessage{},
		lberr.NewErr(int32(lbwebsocket.ErrCode_ErrMessageNotFound), "ErrMessageNotFound"),
	)

	VisitorOrm = NewOrmCondBuilder(
		&lbwebsocket.ModelVisitor{},
		lberr.NewErr(int32(lbwebsocket.ErrCode_ErrVisitorNotFound), "ErrVisitorNotFound"),
	)
}
