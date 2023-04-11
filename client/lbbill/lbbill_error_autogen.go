// Code generated by gen_errorcode.go, DO NOT EDIT.

package lbbill

import (
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	Success             = lberr.NewErr(int32(ErrCode_Success), "Success")
	ErrBillNotFound     = lberr.NewErr(int32(ErrCode_ErrBillNotFound), "ErrBillNotFound")
	ErrCategoryNotFound = lberr.NewErr(int32(ErrCode_ErrCategoryNotFound), "ErrCategoryNotFound")
)