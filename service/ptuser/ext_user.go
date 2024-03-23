package ptuser

import "github.com/oldbai555/bgg/service/lbuser"

func ConvertBaseUser(user *lbuser.ModelUser) *lbuser.BaseUser {
	return &lbuser.BaseUser{
		Id:       user.Id,
		Avatar:   user.Avatar,
		Nickname: user.Nickname,
	}
}
