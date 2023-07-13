package lbddz

import (
	"github.com/oldbai555/bgg/client/lbddz/iface"
)

const (
	ConsumeTypeMatch        iface.ConsumeType = iota + 1 // 匹配对局
	ConsumeTypeStartGame                                 // 开始对局
	ConsumeTypeWantLandlord                              // 叫地主
	ConsumeTypePlayCard                                  // 出牌

	OrmConsumeTypePlayerLogin  // 更新玩家登录
	OrmConsumeTypePlayerLogout // 更新玩家登出
	OrmConsumeTypeLoadPlayer   // 加载玩家数据
	OrmConsumeTypeSyncGameData // 同步对局数据

	WebhookConsumeTypeWebHook // 推送webhook
	ConsumeTypeEnd            // 最后一个消息类型
)
