package service

import (
	"github.com/oldbai555/bgg/client/lbaccount"
)

var AccountServer LbaccountServer

type LbaccountServer struct {
	*lbaccount.UnimplementedLbaccountServer
}
