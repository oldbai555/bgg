package entity

// Game Character

type Actor struct {
	Id        uint64 `json:"id"`
	AccountId uint64 `json:"account_id"`

	Nickname string `json:"nickname"`
	Sex      int    `json:"sex"`
}
