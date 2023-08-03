// Code generated by gen_errorcode.go, DO NOT EDIT.

package lbddz

import (
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	Success               = lberr.NewErr(int32(ErrCode_Success), "Success")
	ErrAlreadyRegister    = lberr.NewErr(int32(ErrCode_ErrAlreadyRegister), "用户名已被注册")
	ErrPasswordMistake    = lberr.NewErr(int32(ErrCode_ErrPasswordMistake), "密码错误")
	ErrPlayerNotFound     = lberr.NewErr(int32(ErrCode_ErrPlayerNotFound), "ErrPlayerNotFound")
	ErrRoomNotFound       = lberr.NewErr(int32(ErrCode_ErrRoomNotFound), "ErrRoomNotFound")
	ErrGameNotFound       = lberr.NewErr(int32(ErrCode_ErrGameNotFound), "ErrGameNotFound")
	ErrGamePlayerNotFound = lberr.NewErr(int32(ErrCode_ErrGamePlayerNotFound), "ErrGamePlayerNotFound")
	ErrPlayCardNotFound   = lberr.NewErr(int32(ErrCode_ErrPlayCardNotFound), "ErrPlayCardNotFound")
)