package actorsys

import "github.com/oldbai555/bgg/service/lbgame/iface"

var _ iface.ISystem = (*Base)(nil)

type Base struct {
	typ   iface.SystemType
	mgr   iface.ISystemMgr
	owner iface.IActor
}

func (b *Base) GetSysType() iface.SystemType {
	return b.typ
}

func (b *Base) GetMgr() iface.ISystemMgr {
	return b.mgr
}

func (b *Base) GetOwner() iface.IActor {
	return b.owner
}

func (b *Base) Init(typ iface.SystemType, mgr iface.ISystemMgr, owner iface.IActor) {
	b.typ = typ
	b.mgr = mgr
	b.owner = owner
}

func (b *Base) OnInit() {
}
