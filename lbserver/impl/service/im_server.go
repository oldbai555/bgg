package service

import (
	"github.com/oldbai555/bgg/client/lbim"
)

var ImServer LbimServer

type LbimServer struct {
	*lbim.UnimplementedLbimServer
}
