// Code generated by gen_errorcode.go, DO NOT EDIT.

package lbwebsocket

import (
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	Success              = lberr.NewErr(int32(ErrCode_Success), "Success")
	ErrCodeConnectClosed = lberr.NewErr(int32(ErrCode_ErrCodeConnectClosed), "ErrCodeConnectClosed")
	ErrChatNotFound      = lberr.NewErr(int32(ErrCode_ErrChatNotFound), "ErrChatNotFound")
	ErrVisitorNotFound   = lberr.NewErr(int32(ErrCode_ErrVisitorNotFound), "ErrVisitorNotFound")
	ErrMessageNotFound   = lberr.NewErr(int32(ErrCode_ErrMessageNotFound), "ErrMessageNotFound")
)