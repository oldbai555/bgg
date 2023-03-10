package impl

import (
	"encoding/xml"
)

type WeChatGzhConf struct {
	AppId          string `json:"app_id"`
	AppSecret      string `json:"app_secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
}

const (
	MsgTypeText     = "text"
	MsgTypeImage    = "image"
	MsgTypeVoice    = "voice"
	MsgTypeVideo    = "video"
	MsgTypeLocation = "location"
	MsgTypeLink     = "link"
	MsgTypeEvent    = "event"
)

const (
	EventTypeSubscribe   = "subscribe"   // 扫描带参数二维码事件 用户未关注时，进行关注后的事件推送
	EventTypeUnsubscribe = "unsubscribe" // 关注/取消关注事件

	EventTypeScan     = "SCAN"     // 扫描带参数二维码事件 用户已关注时的事件推送
	EventTypeLocation = "LOCATION" // 上报地理位置事件
	EventTypeClick    = "CLICK"    // 点击菜单拉取消息时的事件推送
	EventTypeView     = "VIEW"     // 点击菜单跳转链接时的事件推送
)

const (
	SpeechAnswer      = "获取答案"
	SpeechAnswerFail  = "对不起,我找不到您获取的答案，您可以向我提问。例如：\n\n提问 今天是星期几\n\n进行提问，注意中间的空格用于区分。"
	SpeechAnswerReady = "答案还在生成，请稍后再试。"

	SpeechAsk          = "提问"
	SpeechAskSensitive = "对不起，您的提问比较敏感。"
	SpeechAskFail      = "对不起，我找不到您提问的内容，请按示例进行提问。例如：\n\n提问 今天是星期几\n\n进行提问，注意中间的空格用于区分。"
	SpeechErr          = "对不起，我还在学习中。"

	SpeechQueueStartTemplate = "您的排号为：%s,一分钟后重新获取结果\n例如：\n获取答案 xxxxxxxx\n"
)

// CallBackData 微信回调
type CallBackData struct {
	ToUserName   string // 开发者微信号
	FromUserName string // 发送方帐号（一个OpenID）
	CreateTime   int64  // 消息创建时间 （整型）
	MsgType      string // 消息类型
	MsgId        int64  // 消息id，64位整型

	Content string // 文本消息内容

	MediaId string // 消息媒体id，可以调用获取临时素材接口拉取数据。

	PicUrl string // 图片链接（由系统生成）

	Format      string // 语音格式，如amr，speex等
	Recognition string // 开通语音识别后,语音识别结果，UTF8编码

	ThumbMediaId string // 视频消息缩略图的媒体id，可以调用多媒体文件下载接口拉取数据。

	Location_X float64 // 地理位置纬度
	Location_Y float64 // 地理位置经度
	Scale      uint64  // 地图缩放大小
	Label      string  // 地理位置信息

	Title       string // 消息标题
	Description string // 消息描述
	Url         string // 消息链接

	Event    string // 事件类型，subscribe(订阅)、unsubscribe(取消订阅)
	EventKey string // 事件 KEY 值

	Latitude  float64 // 地理位置纬度
	Longitude float64 // 地理位置经度
	Precision float64 // 地理位置精度

	MsgDataId int64 // 消息的数据ID（消息如果来自文章时才有）
	Idx       int64 // 多图文时第几篇文章，从1开始（消息如果来自文章时才有）
}

// WXRepTextMsg 微信回复文本消息结构体
type WXRepTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string

	// 若不标记XMLName, 则解析后的xml名为该结构体的名称
	XMLName xml.Name `xml:"xml"`
}
