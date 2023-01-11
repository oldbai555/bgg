package impl

import (
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/lbtool/pkg/lberr"
)

var (
	ErrPasswordInvalid = lberr.NewErr(int32(lbuser.ErrCode_ErrPasswordInvalid), "密码错误")
	ErrNotLoginInfo    = lberr.NewInvalidArg("缺少登录信息")
	ErrGetLoginFail    = lberr.NewInvalidArg("获取登录信息失败")
)

const (
	RedisPrefix  = "lbblog_"
	CtxWithClaim = "claim"
	LogWithHint  = "hint"
)
