package impl

import (
	"github.com/oldbai555/bgg/client/lbim"
)

var lbimServer LbimServer

type LbimServer struct {
	*lbim.UnimplementedLbimServer
}
