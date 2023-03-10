package impl

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/bgg/lbserver/impl/gpt"
	"github.com/oldbai555/bgg/lbserver/impl/wordscheck"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/routine"
	"github.com/oldbai555/lbtool/utils"
	"strings"
	"time"
)

// 处理回调
func doHandlerWXMsgReceive(callBackData *CallBackData) (*WXRepTextMsg, error) {
	var rsp = WXRepTextMsg{
		ToUserName:   callBackData.FromUserName,
		FromUserName: callBackData.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      MsgTypeText,
	}

	var err error
	switch callBackData.MsgType {
	case MsgTypeText:
		rsp.Content, err = doHandlerMsgTypeText(callBackData)
	case MsgTypeImage:
		rsp.Content = SpeechErr
	case MsgTypeVoice:
		rsp.Content = SpeechErr
	case MsgTypeVideo:
		rsp.Content = SpeechErr
	case MsgTypeLocation:
		rsp.Content = SpeechErr
	case MsgTypeLink:
		rsp.Content = SpeechErr
	case MsgTypeEvent:
		rsp.Content, err = doHandlerWxEvent(callBackData)
	default:
		rsp.Content = SpeechErr
	}

	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

// 处理事件回调
func doHandlerWxEvent(callBackData *CallBackData) (content string, err error) {
	switch callBackData.Event {
	case EventTypeSubscribe:
		content = "谢谢你那么好看还可以来关注我 ~ 我是一个热爱技术的公众号，你可以向我提问，但我不一定会。"
	case EventTypeUnsubscribe:
		content = "希望下一次你还能关注我"
	case EventTypeScan:
	case EventTypeLocation:
	case EventTypeClick:
	case EventTypeView:
	default:
		content = SpeechErr
	}
	return
}

// 处理文本回调
func doHandlerMsgTypeText(callBackData *CallBackData) (string, error) {
	ctx := context.TODO()
	var content = SpeechAnswerTip

	// 检查是否有敏感词
	exit, err := wordscheck.DoCheckWords("http://www.wordscheck.com/wordcheck", callBackData.Content)
	if err != nil {
		log.Errorf("err:%v", err)
		return SpeechErr, err
	}

	if exit {
		return SpeechAskSensitive, nil
	}

	// 检查是否是获取结果
	if strings.HasPrefix(callBackData.Content, SpeechAnswer) {
		split := strings.Split(callBackData.Content, " ")
		if len(split) != 2 {
			return SpeechAnswerFail, nil
		}

		// 拿到等待编号
		var uuid = split[1]

		// 从 redis 拿
		r, err := lb.Rdb.Get(ctx, uuid).Result()
		if err != nil && err != redis.Nil {
			log.Errorf("err:%v", err)
			return SpeechAnswerFail, nil
		}

		// 结果为空
		if err == redis.Nil {
			return SpeechAnswerReady, nil
		}

		// 找到啦
		if err == nil {
			err := lb.Rdb.Del(ctx, uuid).Err()
			if err != nil {
				log.Errorf("err:%v", err)
			}
			return r, nil
		}

		// 应该不会走到这里吧
		return content, nil
	}

	// 检查是否是提问
	if strings.HasPrefix(callBackData.Content, SpeechAsk) {

		split := strings.Split(callBackData.Content, " ")
		if len(split) != 2 {
			return SpeechAskFail, nil
		}

		// 生成等待编号
		uuid := utils.GenUUID()

		// 异步去获取结果
		routine.Go(ctx, func(ctx context.Context) error {
			// 去查gpt
			completionRsp, err := gpt.DoChatCompletion(lb.V.GetString("chatGpt.proxy"), lb.V.GetString("chatGpt.api_key"), &gpt.ChatCompletionReq{
				Model: gpt.DefualtModel,
				Messages: []*gpt.ChatCompletionReqMessage{
					{
						Content: split[1],
						Role:    gpt.DefualtRole,
					},
				},
			})

			// 出错啦
			if err != nil {
				log.Errorf("err:%v", err) // 一般是请求不过去或者网络超时
				err = SetGptResult(ctx, uuid, SpeechErr)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				return err
			}

			// 写入结果
			var result string
			for _, choice := range completionRsp.Choices {
				if choice.Message != nil {
					result = result + choice.Message.Content
				}
			}

			// 兜底回复
			if len(result) == 0 {
				result = SpeechErr
			}

			if len(result) > 0 {
				if strings.Contains(result, "AI语言模型") {
					strings.ReplaceAll(result, "AI语言模型", "公众号")
				}
				if strings.Contains(result, "AI") {
					strings.ReplaceAll(result, "AI", "公众号")
				}
				if strings.Contains(result, "公众号公众号") {
					strings.ReplaceAll(result, "公众号公众号", "公众号")
				}
			}

			// 写入redis
			err = SetGptResult(ctx, uuid, result)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		})

		// 等三秒后看看能不能直接拿到结果
		timer := time.NewTimer(3 * time.Second)
		select {
		case <-timer.C:
			content, err = GetGptResult(ctx, uuid)
			if err != nil {
				log.Errorf("err:%v", err)
				return content, err
			}
		}

		// 对接收的消息进行被动回复
		return content, nil
	}

	// 返回结果
	return content, nil
}

// GetGptResult 获取gpt的结果
func GetGptResult(ctx context.Context, uuid string) (content string, err error) {
	content, err = lb.Rdb.Get(ctx, uuid).Result()
	if err != nil && err != redis.Nil {
		log.Errorf("err:%v", err)
		return
	}
	// 没错就结束咯
	if err == nil {
		err = lb.Rdb.Del(ctx, uuid).Err()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		return
	}
	// 拿不到就给个默认的话语
	return fmt.Sprintf(SpeechQueueStartTemplate, uuid), nil
}

// SetGptResult 写入gpt的结果
func SetGptResult(ctx context.Context, uuid string, content string) error {
	err := lb.Rdb.Set(ctx, uuid, content, time.Hour).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
