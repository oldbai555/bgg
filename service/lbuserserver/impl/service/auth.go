/**
 * @Author: zjj
 * @Date: 2024/3/23
 * @Desc:
**/

package service

import (
	"context"
	emoji "github.com/go-xman/go.emoji"
	"github.com/oldbai555/bgg/service/lbuser"
	"github.com/oldbai555/lbtool/log"
	"unicode"
	"unicode/utf8"
)

func (a *LbuserServer) CheckLoginUser(ctx context.Context, sid string) (interface{}, error) {
	userSysRsp, err := a.GetLoginUser(ctx, &lbuser.GetLoginUserReq{
		Sid: sid,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return userSysRsp.BaseUser, nil
}

func checkPassword(pwd string) bool {
	if !CheckStr(pwd) {
		return false
	}
	if !CheckSpecialCharacter(pwd) {
		return false
	}
	return true
}

func CheckStr(str string) bool {
	if len(str) == 0 {
		return false
	}
	//名字最大不能超过 16 个字符
	pwdLen := utf8.RuneCountInString(str)
	if pwdLen > 16 {
		return false
	}
	return true
}

func CheckSpecialCharacter(name string) bool {
	// 空白符号
	for _, c := range name {
		if unicode.IsSpace(c) {
			return false
		}
	}
	// emoji 表情
	if emoji.HasEmoji(name) {
		return false
	}
	return true
}
