package impl

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"sort"
)

func registerWechatGzhApi(h *gin.Engine) {
	group := h.Group("wx")

	// 可以利用反射来映射函数进去
	group.GET("/wechatgzh", WXGzhAuth)
	group.POST("/wechatgzh", WXMsgReceive)
}

func WXGzhAuth(c *gin.Context) {
	handler := NewHandler(c)

	signature, _ := handler.GetQuery("signature")
	timestamp, _ := handler.GetQuery("timestamp")
	echostr, _ := handler.GetQuery("echostr")
	nonce, _ := handler.GetQuery("nonce")

	// 3.token，timestamp，nonce按字典排序的字符串list
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

// WXMsgReceive 微信消息接收
func WXMsgReceive(c *gin.Context) {
	var callBackData CallBackData
	err := c.ShouldBindXML(&callBackData)
	if err != nil {
		log.Infof("[消息接收] - XML数据包解析失败: %v", err)
		return
	}

	log.Infof("[消息接收] - 收到消息, 消息类型为: %s , 消息内容: %v", callBackData.MsgType, callBackData)
	reply, err := doHandlerWXMsgReceive(&callBackData)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	WXMsgReply(c, reply)
}

// WXMsgReply 微信消息回复
func WXMsgReply(c *gin.Context, reply *WXRepTextMsg) {
	msg, err := xml.Marshal(&reply)
	if err != nil {
		log.Infof("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}
