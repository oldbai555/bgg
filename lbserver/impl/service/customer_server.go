package service

import (
	"github.com/oldbai555/bgg/client/lbcustomer"
)

var CustomerServer LbcustomerServer

type LbcustomerServer struct {
	*lbcustomer.UnimplementedLbcustomerServer
}
