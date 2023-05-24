// Code generated by gen_errorcode.go, DO NOT EDIT.

package lbblog

import (
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	Nil                 = lberr.NewErr(int32(ErrCode_Nil), "Nil")
	ErrArticleNotFound  = lberr.NewErr(int32(ErrCode_ErrArticleNotFound), "ErrArticleNotFound")
	ErrCategoryNotFound = lberr.NewErr(int32(ErrCode_ErrCategoryNotFound), "ErrCategoryNotFound")
	ErrCommentNotFound  = lberr.NewErr(int32(ErrCode_ErrCommentNotFound), "ErrCommentNotFound")
)
