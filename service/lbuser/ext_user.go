package lbuser

func NewBaseUser(u *ModelUser) *BaseUser {
	return &BaseUser{
		Id:       u.Id,
		Username: u.Username,
		Avatar:   u.Avatar,
		Nickname: u.Nickname,
		Email:    u.Email,
		Github:   u.Github,
		Desc:     u.Desc,
		Role:     u.Role,
	}
}
