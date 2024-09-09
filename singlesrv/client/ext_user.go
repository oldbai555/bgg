/**
 * @Author: zjj
 * @Date: 2024/3/25
 * @Desc:
**/

package client

func (x *ModelUser) ToBaseUser() *BaseUser {
	return &BaseUser{
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
