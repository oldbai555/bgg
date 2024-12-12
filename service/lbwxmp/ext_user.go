/**
 * @Author: zjj
 * @Date: 2024/12/12
 * @Desc:
**/

package lbwxmp

import "github.com/oldbai555/bgg/service/lbbase"

func (x *ModelMpMemberUser) ToBaseUser() *lbbase.BaseUser {
	return &lbbase.BaseUser{
		Id:        x.Id,
		Username:  x.Mobile,
		Avatar:    x.Avatar,
		Nickname:  x.Nickname,
		CreatedAt: x.CreatedAt,
		Type:      1,
	}
}
