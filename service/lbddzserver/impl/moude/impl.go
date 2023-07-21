package moude

import (
	"github.com/name5566/leaf/module"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/db"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/game"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/gate"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/login"
)

func Modules() []module.Module {
	return []module.Module{
		db.Module,
		game.Module,
		gate.Module,
		login.Module,
	}
}
