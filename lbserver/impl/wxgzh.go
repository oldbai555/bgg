package impl

import (
	"context"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/bgg/lbserver/impl/gpt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"sort"
	"strings"
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

	var result string
	if strings.HasPrefix(textMsg.Content, "获取答案:") {
		split := strings.Split(textMsg.Content, ":")
		if len(split) != 2 {
			WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, "对不起,我找不到你想要的答案,请按格式=>\n获取答案:xxxxxxxx\n获取结果。")
			return
		}
		r, err := lb.Rdb.Get(c, split[1]).Result()
		if err != nil && err != redis.Nil {
			log.Errorf("err:%v", err)
			return
		}
		if err == redis.Nil {
			result = "答案还在生成中,请稍等"
		}
		if err == nil {
			result = r
		}
		WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, result)
		return
	}

	if strings.HasPrefix("提问:", textMsg.Content) {
		split := strings.Split(textMsg.Content, ":")
		if len(split) != 2 {
			WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, "对不起,我找不到你想要的提问的内容,请按格式=>\n提问:今天是星期几\n进行提问")
			return
		}
		uuid := utils.GenUUID()
		// 异步去处理
		routine.Go(c, func(ctx context.Context) error {
			completionRsp, err := gpt.DoChatCompletion(lb.V.GetString("chatGpt.api_key"), &gpt.ChatCompletionReq{
				Model: "gpt-3.5-turbo",
				Messages: []*gpt.ChatcompletionreqMessage{
					{
						Content: textMsg.Content,
						Role:    "user",
					},
				},
			})
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			var result string
			for _, choice := range completionRsp.Choices {
				if choice.Message != nil {
					result = result + choice.Message.Content
				}
			}

			if len(result) == 0 {
				result = "对不起,你说的我暂时不能理解"
			}
			err = lb.Rdb.Set(ctx, uuid, result, time.Hour*24).Err()
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			return nil
		})

		timer := time.NewTimer(3 * time.Second)
		select {
		case <-timer.C:
			r, err := lb.Rdb.Get(c, uuid).Result()
			if err != nil && err != redis.Nil {
				log.Errorf("err:%v", err)
				return
			}
			if err == redis.Nil {
				result = fmt.Sprintf("你的等待序号为:%s,一分钟后请按格式=>\n获取答案:xxxxxxxx\n获取结果", uuid)
			}
			if err == nil {
				result = r
			}
		}

		// 对接收的消息进行被动回复
		WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, result)
		return
	}

	WXMsgReply(c, textMsg.ToUserName, textMsg.FromUserName, "请按格式=>\n提问:今天是星期几\n进行提问")
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
func WXMsgReply(c *gin.Context, fromUser, toUser, Content string) {
	repTextMsg := WXRepTextMsg{
		ToUserName:   toUser,
		FromUserName: fromUser,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      Content,
	}

	msg, err := xml.Marshal(&repTextMsg)
	if err != nil {
		log.Infof("[消息回复] - 将对象进行XML编码出错: %v\n", err)
		return
	}
	_, _ = c.Writer.Write(msg)
}
