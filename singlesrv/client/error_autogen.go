// Code generated by gen_error_code.go, DO NOT EDIT.

package client

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	Success                      = status.Error(codes.Code(int32(ErrCode_Success)), "Success")
	ErrFileNotFound              = status.Error(codes.Code(int32(ErrCode_ErrFileNotFound)), "文件不存在")
	ErrUserNotFound              = status.Error(codes.Code(int32(ErrCode_ErrUserNotFound)), "用户不存在")
	ErrOldPasswordNotEqual       = status.Error(codes.Code(int32(ErrCode_ErrOldPasswordNotEqual)), "旧密码有误")
	ErrOldPwdEqualNewPwd         = status.Error(codes.Code(int32(ErrCode_ErrOldPwdEqualNewPwd)), "新旧密码不能相同")
	ErrNsqProducerConnectFailure = status.Error(codes.Code(int32(ErrCode_ErrNsqProducerConnectFailure)), "连接nsq生产者失败")
	ErrNsqTopicAlready           = status.Error(codes.Code(int32(ErrCode_ErrNsqTopicAlready)), "nsq topic 已经存在")
	ErrFileMd5IsEmpty            = status.Error(codes.Code(int32(ErrCode_ErrFileMd5IsEmpty)), "文件 md5 不能为空")
	ErrFileMd5Already            = status.Error(codes.Code(int32(ErrCode_ErrFileMd5Already)), "文件 md5 已经存在")
	ErrFileUploadFailure         = status.Error(codes.Code(int32(ErrCode_ErrFileUploadFailure)), "文件上传失败")
	ErrFileAlreadyExist          = status.Error(codes.Code(int32(ErrCode_ErrFileAlreadyExist)), "保存文件失败，文件重复")
)