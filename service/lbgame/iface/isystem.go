package iface

type SystemType uint32

type GenISystem func() ISystem

type ISystemMgr interface{}

type ISystem interface {
	Init(typ SystemType, mgr ISystemMgr, owner IActor)
	OnInit()

	GetSysType() SystemType
	GetMgr() ISystemMgr
	GetOwner() IActor
}
