package engine

import "github.com/oldbai555/bgg/service/lbgame/iface"

// SysSet 全局的角色系统构造函数
var SysSet = make(map[iface.SystemType]iface.GenISystem)
