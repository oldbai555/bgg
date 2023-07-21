package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/conf"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/moude/game"
	"github.com/oldbai555/bgg/service/lbddzserver/impl/msg"
)

type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.Server.WSAddr,
		HTTPTimeout:     conf.HTTPTimeout,
		CertFile:        conf.Server.CertFile,
		KeyFile:         conf.Server.KeyFile,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       msg.Processor, // 指向全局的processor
		AgentChanRPC:    game.ChanRPC,  // 指向游戏模块的AgentChanRPC
	}
}
