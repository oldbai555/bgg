package lbwxmpserver

import (
	"github.com/oldbai555/bgg/service/lbwxmp"
)

var OnceSvrImpl = &LbwxmpServer{}

type LbwxmpServer struct {
	lbwxmp.UnimplementedLbwxmpServer
}
