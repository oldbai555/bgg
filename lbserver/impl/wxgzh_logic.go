package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/bgg/client/lbaccount"
	"github.com/oldbai555/bgg/client/lbcustomer"
	"github.com/oldbai555/bgg/client/lbim"
	"github.com/oldbai555/bgg/lbserver/impl/gpt"
	"github.com/oldbai555/bgg/lbserver/impl/wordscheck"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/restysdk"
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

	var msg = lbim.ModelMessage{
		ServerMsgId: fmt.Sprintf("%d", callBackData.MsgId),
		SendAt:      uint64(callBackData.CreateTime),
		From:        callBackData.FromUserName,
		To:          callBackData.ToUserName,
		Source:      uint32(lbim.MessageSource_MessageSourceWxGzh),
		Content:     &lbim.Content{},
	}

	var err error
	switch callBackData.MsgType {
	case MsgTypeText:
		rsp.Content, err = doHandlerMsgTypeText(callBackData)
		msg.Content.Text = &lbim.Content_Text{
			Content: callBackData.Content,
		}
	case MsgTypeImage:
		rsp.Content = SpeechErr
		url, err := ConvertMediaUrl(fmt.Sprintf("%s.jpg", utils.GenUUID()), callBackData.PicUrl)
		if err != nil {
			log.Errorf("err is %v", err)
			url = callBackData.PicUrl
		}
		msg.Content.Image = &lbim.Content_Image{
			Url: url,
		}
	case MsgTypeVoice:
		// 语音获取资源是得到一个字节流文件，可以直接下载
		rsp.Content = SpeechErr
		bytes, err := GetWxGzhMediaBytes(context.TODO(), callBackData.MediaId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		url, err := ConvertMediaBytes(fmt.Sprintf("%s.%s", utils.GenUUID(), callBackData.Format), bytes)
		if err != nil {
			log.Errorf("err:%v", err)
			url = callBackData.MediaId
		}
		msg.Content.Voice = &lbim.Content_Voice{
			Url:         url,
			Format:      callBackData.Format,
			Recognition: callBackData.Recognition,
		}
	case MsgTypeVideo:
		rsp.Content = SpeechErr
		bytes, err := GetWxGzhMediaBytes(context.TODO(), callBackData.MediaId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		url, err := ConvertMediaBytes(fmt.Sprintf("%s.mp4", utils.GenUUID()), bytes)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		bytes, err = GetWxGzhMediaBytes(context.TODO(), callBackData.ThumbMediaId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		videoPicUrl, err := ConvertMediaBytes(fmt.Sprintf("%s.jpg", utils.GenUUID()), bytes)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		msg.Content.Video = &lbim.Content_Video{
			Url:     url,
			Caption: videoPicUrl, // todo 应有一个字段来承接视频封面
		}
	case MsgTypeLocation:
		rsp.Content = SpeechErr
		msg.Content.Location = &lbim.Content_Location{
			X:     callBackData.Location_X,
			Y:     callBackData.Location_Y,
			Scale: float64(callBackData.Scale),
			Label: callBackData.Label,
		}
	case MsgTypeLink:
		rsp.Content = SpeechErr
		msg.Content.Document = &lbim.Content_Document{
			Url:     callBackData.Url,
			Caption: fmt.Sprintf("标题：%s,描述：%s", callBackData.Title, callBackData.Description),
		}
	case MsgTypeEvent:
		rsp.Content, err = doHandlerWxEvent(callBackData)
	default:
		rsp.Content = SpeechErr
	}

	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 插入消息记录
	err = MessageOrm.NewScope().Create(context.TODO(), &msg)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 插入客户
	err = CustomerOrm.NewScope().UpdateOrCreate(context.TODO(), map[string]interface{}{
		lbcustomer.FieldSn_:      callBackData.FromUserName,
		lbcustomer.FieldChannel_: uint32(lbcustomer.Channel_ChannelWx),
	}, map[string]interface{}{
		lbcustomer.FieldSn_:      callBackData.FromUserName,
		lbcustomer.FieldChannel_: uint32(lbcustomer.Channel_ChannelWx),
	}, &lbcustomer.ModelCustomer{})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 插入系统账号
	err = AccountOrm.NewScope().UpdateOrCreate(context.TODO(), map[string]interface{}{
		lbaccount.FieldSn_:      callBackData.ToUserName,
		lbaccount.FieldChannel_: uint32(lbaccount.Channel_ChannelWxGzh),
	}, map[string]interface{}{
		lbaccount.FieldSn_:      callBackData.ToUserName,
		lbaccount.FieldChannel_: uint32(lbaccount.Channel_ChannelWxGzh),
	}, &lbaccount.ModelAccount{})
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

// GetWxGzhAccessToken 获取 accessToken
func GetWxGzhAccessToken(ctx context.Context) (string, error) {
	accessToken, err := lb.Rdb.Get(ctx, fmt.Sprintf("%s_%s", lb.WechatConf.AppId, lb.WechatConf.AppSecret)).Result()

	if err != nil && err != redis.Nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	if err == nil && accessToken != "" {
		return accessToken, nil
	}

	// 找不到 那就去微信拿
	resp, err := restysdk.NewRequest().SetQueryParams(map[string]string{
		"grant_type": "client_credential",
		"appid":      lb.WechatConf.AppId,
		"secret":     lb.WechatConf.AppSecret,
	}).Get(WxGzhAccessTokenUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	log.Debugf("url is %s", resp.Request.URL)

	if err = CheckWxGzhApiErr(resp.Body()); err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	var accessResp AccessTokenResp
	err = json.Unmarshal(resp.Body(), &accessResp)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	if accessResp.AccessToken == "" {
		log.Errorf("accessResp.AccessToken is null")
		return "", lberr.NewInvalidArg("accessResp.AccessToken is null")
	}

	// 写入redis
	if accessResp.ExpiresIn == 0 {
		accessResp.ExpiresIn = 7200
	}

	err = lb.Rdb.Set(ctx, fmt.Sprintf("%s_%s", lb.WechatConf.AppId, lb.WechatConf.AppSecret), accessResp.AccessToken, time.Duration(accessResp.ExpiresIn)*time.Second).Err()
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	return "", nil
}

// GetWxGzhMediaBytes 获取临时媒体资源的字节流
func GetWxGzhMediaBytes(ctx context.Context, mediaId string) ([]byte, error) {
	token, err := GetWxGzhAccessToken(ctx)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	resp, err := restysdk.NewRequest().SetQueryParams(map[string]string{
		"access_token": token,
		"media_id":     mediaId,
	}).Get(WxGzhMediaUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return resp.Body(), nil
}

// CheckWxGzhApiErr 检查是否错误
func CheckWxGzhApiErr(data []byte) error {
	var apiErr WxGzhApiErr
	err := json.Unmarshal(data, &apiErr)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if apiErr.ErrCode != 0 {
		return &apiErr
	}

	return nil
}
