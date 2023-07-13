package cc

import (
	"errors"
	"github.com/oldbai555/baix/iface"
	"github.com/oldbai555/lbtool/log"
	"sync"
)

type ConnIDMap map[uint64]iface.IConnection

type notify struct {
	cimap ConnIDMap
	sync.RWMutex
}

var Notify = NewNotify()

func NewNotify() *notify {
	return &notify{
		cimap: make(map[uint64]iface.IConnection, 5000),
	}
}

func (n *notify) ConnNums() int {
	return len(n.cimap)
}

func (n *notify) HasIdConn(Id uint64) bool {
	n.RLock()
	defer n.RUnlock()
	_, ok := n.cimap[Id]
	return ok
}

func (n *notify) SetNotifyID(Id uint64, conn iface.IConnection) {
	n.Lock()
	defer n.Unlock()
	n.cimap[Id] = conn
}

func (n *notify) GetNotifyByID(Id uint64) (iface.IConnection, error) {
	n.RLock()
	defer n.RUnlock()
	Conn, ok := n.cimap[Id]
	if !ok {
		return nil, errors.New(" Not Find UserId")
	}
	return Conn, nil
}

func (n *notify) DelNotifyByID(Id uint64) {
	n.RLock()
	defer n.RUnlock()
	delete(n.cimap, Id)
}

func (n *notify) NotifyBuffToConnByIDList(msgId uint32, data []byte, connIdList ...uint64) error {
	for _, connId := range connIdList {
		Conn, err := n.GetNotifyByID(connId)
		if err != nil {
			return err
		}
		err = Conn.SendBuffMsg(msgId, data)
		if err != nil {
			log.Errorf("Notify to %d err:%s \n", connId, err)
			return err
		}
	}
	return nil
}

func (n *notify) NotifyBuffAll(MsgId uint32, data []byte) error {
	n.RLock()
	defer n.RUnlock()
	for Id, v := range n.cimap {
		err := v.SendBuffMsg(MsgId, data)
		if err != nil {
			log.Errorf("Notify to %d err:%s \n", Id, err)
		}
	}
	return nil
}
