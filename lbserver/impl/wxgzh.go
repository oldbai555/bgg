package impl

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"sort"
	"time"
)

type WeChatGzhConf struct {
	AppId          string `json:"app_id"`
	AppSecret      string `json:"app_secret"`
	Token          string `json:"token"`
	EncodingAESKey string `json:"encoding_aes_key"`
}

func registerWechatGzhApi(h *gin.Engine) {
	group := h.Group("wx")

	// 可以利用反射来映射函数进去
	group.GET("/wechatgzh", AuthWeChatGzh)
	group.POST("/wechatgzh", WXMsgReceive)
}

func AuthWeChatGzh(c *gin.Context) {
	handler := NewHandler(c)

	signature, _ := handler.GetQuery("signature")
	timestamp, _ := handler.GetQuery("timestamp")
	echostr, _ := handler.GetQuery("echostr")
	nonce, _ := handler.GetQuery("nonce")

	//3.token，timestamp，nonce按字典排序的字符串list
	strs := sort.StringSlice{lb.WechatConf.Token, timestamp, nonce} // 使用本地的token生成校验
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}

	// 4. 哈希算法加密list得到hashcode
	h := sha1.New()
	h.Write([]byte(str))
	hashcode := fmt.Sprintf("%x", h.Sum(nil)) // h.Sum(nil) 做hash

	var rspStr = echostr
	// 比对数据
	if signature != hashcode {
		rspStr = "error"
	}

	// 写入结果
	c.Header(HttpHeaderContentType, "text/plain;charset=utf-8")
	_, err := c.Writer.WriteString(rspStr)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	return
}

// WXTextMsg 微信文本消息结构体
type WXTextMsg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	Content      string
	MsgId        int64
}

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	var textMsg WXTextMsg
	err := c.ShouldBindXML(&textMsg)
	if err != nil {
		log.Infof("[消息接收] - XML数据包解析失败: %v", err)
		return
	}

	log.Infof("[消息接收] - 收到消息, 消息类型为: %s, 消息内容为: %s", textMsg.MsgType, textMsg.Content)

	// 对接收的消息进行被动回复
	WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName)
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

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, fromUser, toUser string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      fmt.Sprintf("[消息回复] - %s", time.Now().Format("2006-01-02 15:04:05")),
	}

	msg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Infof("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}
