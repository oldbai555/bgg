package impl

import (
	"crypto/sha1"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"sort"
	"strings"
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
	group.POST("/wechatgzh", DoWeChatCallBack)
}

func AuthWeChatGzh(c *gin.Context) {
	handler := NewHandler(c)

	signature, _ := handler.GetQuery("signature")
	timestamp, _ := handler.GetQuery("timestamp")
	echostr, _ := handler.GetQuery("echostr")
	nonce, _ := handler.GetQuery("nonce")

	strList := []string{lb.WechatConf.Token, timestamp, nonce}
	// 字典排序
	sort.Strings(strList)

	// sha1 加密
	h := sha1.New()
	sum := h.Sum([]byte(strings.Join(strList, "")))

	var rspStr = echostr
	// 比对数据
	if signature != string(sum) {
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

func DoWeChatCallBack(c *gin.Context) {

	handler := NewHandler(c)

	signature, _ := handler.GetQuery("signature")
	timestamp, _ := handler.GetQuery("timestamp")
	echostr, _ := handler.GetQuery("echostr")
	nonce, _ := handler.GetQuery("nonce")

	strList := []string{lb.WechatConf.Token, timestamp, nonce}
	// 字典排序
	sort.Strings(strList)

	// sha1 加密
	h := sha1.New()
	sum := h.Sum([]byte(strings.Join(strList, "")))

	var rspStr = echostr
	// 比对数据
	if signature != string(sum) {
		rspStr = "error"
	}

	// 写入结果
	c.Header(HttpHeaderContentType, "text/plain;charset=utf-8")
	_, err := c.Writer.WriteString(rspStr)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

}
