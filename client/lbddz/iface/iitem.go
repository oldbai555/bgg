package iface

type IItem interface {
	GetMsgType() ConsumeType
	GetParams() []interface{}
}
