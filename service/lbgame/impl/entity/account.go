package entity

type Account struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
