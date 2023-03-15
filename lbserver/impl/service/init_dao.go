package service

import (
	"context"
	"github.com/oldbai555/bgg/lbserver/impl/dao"
	"github.com/oldbai555/bgg/lbserver/impl/dao/impl/mysql"
	"github.com/oldbai555/lbtool/log"
)

var (
	Account  dao.AccountDao
	Article  dao.ArticleDao
	Category dao.CategoryDao
	Comment  dao.CommentDao
	Customer dao.CustomerDao
	Message  dao.MessageDao
	User     dao.UserDao
)

func InitDao(ctx context.Context, dsn string) (err error) {
	Account, err = mysql.NewAccountImpl(ctx, dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	Article, err = mysql.NewArticleImpl(ctx, dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	Category, err = mysql.NewCategoryImpl(ctx, dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	Comment, err = mysql.NewCommentImpl(ctx, dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	Customer, err = mysql.NewCustomerImpl(ctx, dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	Message, err = mysql.NewMessageImpl(ctx, dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	User, err = mysql.NewUserImpl(ctx, dsn)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	return
}
