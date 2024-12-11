/**
 * @Author: zjj
 * @Date: 2024/3/25
 * @Desc:
**/

package lbsingle

import (
	"github.com/oldbai555/bgg/service/lbbase"
)

func (x *ModelUser) ToBaseUser() *lbbase.BaseUser {
	return &lbbase.BaseUser{
		Id:        x.Id,
		Username:  x.Username,
		Avatar:    x.Avatar,
		Nickname:  x.Nickname,
		Email:     x.Email,
		Github:    x.Github,
		Desc:      x.Desc,
		Role:      x.Role,
		CreatedAt: x.CreatedAt,
	}
}
