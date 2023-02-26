package impl

import (
	"github.com/oldbai555/bgg/client/lbuser"
)

var (
	UserOrm *OrmCondBuilder
)

func InitDbOrm() {
	UserOrm = NewOrmCondBuilder(
		&lbuser.ModelUser{},
		lbuser.ErrUserNotFound,
	)
}
