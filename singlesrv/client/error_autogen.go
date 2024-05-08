// Code generated by gen_error_code.go, DO NOT EDIT.

package client

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Success         = status.Error(codes.Code(int32(ErrCode_Success)), "Success")
	ErrFileNotFound = status.Error(codes.Code(int32(ErrCode_ErrFileNotFound)), "文件不存在")
)
