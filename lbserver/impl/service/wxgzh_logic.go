package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/oldbai555/bgg/client/lbaccount"
	"github.com/oldbai555/bgg/client/lbcustomer"
	"github.com/oldbai555/bgg/client/lbim"
	"github.com/oldbai555/bgg/lbserver/impl/cache"
	"github.com/oldbai555/bgg/lbserver/impl/conf"
	"github.com/oldbai555/bgg/lbserver/impl/constant"
	"github.com/oldbai555/bgg/lbserver/impl/gpt"
	"github.com/oldbai555/bgg/lbserver/impl/storage"
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
func DoHandlerWXMsgReceive(callBackData *constant.CallBackData) (*constant.WXRepTextMsg, error) {
	var rsp = constant.WXRepTextMsg{
		ToUserName:   callBackData.FromUserName,
		FromUserName: callBackData.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      constant.MsgTypeText,
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
	case constant.MsgTypeText:
		rsp.Content, err = doHandlerMsgTypeText(callBackData)
		msg.Content.Text = &lbim.Content_Text{
			Content: callBackData.Content,
		}
	case constant.MsgTypeImage:
		rsp.Content = constant.SpeechErr
		url, err := storage.ConvertMediaUrl(fmt.Sprintf("%s.jpg", utils.GenUUID()), callBackData.PicUrl)
		if err != nil {
			log.Errorf("err is %v", err)
			url = callBackData.PicUrl
		}
		msg.Content.Image = &lbim.Content_Image{
			Url: url,
		}
	case constant.MsgTypeVoice:
		// 语音获取资源是得到一个字节流文件，可以直接下载
		rsp.Content = constant.SpeechErr
		bytes, err := GetWxGzhMediaBytes(context.TODO(), callBackData.MediaId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		url, err := storage.ConvertMediaBytes(fmt.Sprintf("%s.%s", utils.GenUUID(), callBackData.Format), bytes)
		if err != nil {
			log.Errorf("err:%v", err)
			url = callBackData.MediaId
		}
		msg.Content.Voice = &lbim.Content_Voice{
			Url:         url,
			Format:      callBackData.Format,
			Recognition: callBackData.Recognition,
		}
	case constant.MsgTypeVideo:
		rsp.Content = constant.SpeechErr
		bytes, err := GetWxGzhMediaBytes(context.TODO(), callBackData.MediaId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		url, err := storage.ConvertMediaBytes(fmt.Sprintf("%s.mp4", utils.GenUUID()), bytes)
		if err != nil {
			log.Errorf("err:%v", err)
			url = callBackData.MediaId
		}
		bytes, err = GetWxGzhMediaBytes(context.TODO(), callBackData.ThumbMediaId)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		videoPicUrl, err := storage.ConvertMediaBytes(fmt.Sprintf("%s.jpg", utils.GenUUID()), bytes)
		if err != nil {
			log.Errorf("err:%v", err)
			videoPicUrl = callBackData.ThumbMediaId
		}
		msg.Content.Video = &lbim.Content_Video{
			Url:     url,
			Caption: videoPicUrl, // todo 应有一个字段来承接视频封面
		}
	case constant.MsgTypeLocation:
		rsp.Content = constant.SpeechErr
		msg.Content.Location = &lbim.Content_Location{
			X:     callBackData.Location_X,
			Y:     callBackData.Location_Y,
			Scale: float64(callBackData.Scale),
			Label: callBackData.Label,
		}
	case constant.MsgTypeLink:
		rsp.Content = constant.SpeechErr
		msg.Content.Document = &lbim.Content_Document{
			Url:     callBackData.Url,
			Caption: fmt.Sprintf("标题：%s,描述：%s", callBackData.Title, callBackData.Description),
		}
	case constant.MsgTypeEvent:
		rsp.Content, err = doHandlerWxEvent(callBackData)
	default:
		rsp.Content = constant.SpeechErr
	}

	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 插入消息记录
	err = Message.Create(context.TODO(), &msg)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	// 插入客户
	err = Customer.UpdateOrCreate(context.TODO(), map[string]interface{}{
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
	err = Account.UpdateOrCreate(context.TODO(), map[string]interface{}{
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

	// 写入回复消息
	// 将消息写入数据库
	// todo 应该补充消息发送状态
	err = Message.Create(context.TODO(), &lbim.ModelMessage{
		SysMsgId: utils.GenUUID(),
		SendAt:   uint64(rsp.CreateTime),
		From:     rsp.FromUserName,
		To:       rsp.ToUserName,
		Source:   uint32(lbim.MessageSource_MessageSourceWxGzh),
		Content: &lbim.Content{
			Text: &lbim.Content_Text{
				Content: rsp.Content,
			},
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

// 处理事件回调
func doHandlerWxEvent(callBackData *constant.CallBackData) (content string, err error) {
	switch callBackData.Event {
	case constant.EventTypeSubscribe:
		content = "谢谢你那么好看还可以来关注我 ~ 我是一个热爱技术的公众号，你可以向我提问，但我不一定会。"
	case constant.EventTypeUnsubscribe:
		content = "希望下一次你还能关注我"
	case constant.EventTypeScan:
	case constant.EventTypeLocation:
	case constant.EventTypeClick:
	case constant.EventTypeView:
	default:
		content = constant.SpeechErr
	}
	return
}

// 处理文本回调
func doHandlerMsgTypeText(callBackData *constant.CallBackData) (string, error) {
	ctx := context.TODO()
	var content = constant.SpeechAnswerTip

	// 检查是否有敏感词
	exit, err := wordscheck.DoCheckWords("http://www.wordscheck.com/wordcheck", callBackData.Content)
	if err != nil {
		log.Errorf("err:%v", err)
		return constant.SpeechErr, err
	}

	if exit {
		return constant.SpeechAskSensitive, nil
	}

	// 检查是否是获取结果
	if strings.HasPrefix(callBackData.Content, constant.SpeechAnswer) {
		split := strings.Split(callBackData.Content, " ")
		if len(split) != 2 {
			return constant.SpeechAnswerFail, nil
		}

		// 拿到等待编号
		var uuid = split[1]

		// 从 redis 拿
		content, err = cache.GetGptResult(ctx, uuid)
		if err != nil && err != redis.Nil {
			log.Errorf("err:%v", err)
			return constant.SpeechAnswerFail, nil
		}

		// 结果为空
		if err == redis.Nil {
			return constant.SpeechAnswerReady, nil
		}

		// 应该不会走到这里吧
		return content, nil
	}

	// 检查是否是提问
	if strings.HasPrefix(callBackData.Content, constant.SpeechAsk) {

		split := strings.Split(callBackData.Content, " ")
		if len(split) != 2 {
			return constant.SpeechAskFail, nil
		}

		// 生成等待编号
		uuid := utils.GenUUID()

		// 异步去获取结果
		routine.Go(ctx, func(ctx context.Context) error {
			// 去查gpt
			completionRsp, err := gpt.DoChatCompletion(conf.Global.ChatGpt.Proxy, conf.Global.ChatGpt.ApiKey, &gpt.ChatCompletionReq{
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
				err = cache.SetGptResult(ctx, uuid, constant.SpeechErr)
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
				result = constant.SpeechErr
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
			err = cache.SetGptResult(ctx, uuid, result)
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
			content, err = cache.GetGptResult(ctx, uuid)
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

// GetWxGzhAccessToken 获取 accessToken
func GetWxGzhAccessToken(ctx context.Context, appId, appSecret string) (string, error) {
	accessToken, err := cache.GetWxGzhAccessToken(ctx, fmt.Sprintf("%s_%s", appId, appSecret))

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
		"appid":      appId,
		"secret":     appSecret,
	}).Get(constant.WxGzhAccessTokenUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	log.Debugf("url is %s", resp.Request.URL)

	if err = CheckWxGzhApiErr(resp.Body()); err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	var accessResp constant.AccessTokenResp
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

	err = cache.SetWxGzhAccessToken(ctx, fmt.Sprintf("%s_%s", appId, appSecret), accessResp.AccessToken, time.Duration(accessResp.ExpiresIn)*time.Second)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	return accessResp.AccessToken, nil
}

// GetWxGzhMediaBytes 获取临时媒体资源的字节流
func GetWxGzhMediaBytes(ctx context.Context, mediaId string) ([]byte, error) {
	var appId, appSecret = conf.Global.WxGzhConf.AppId, conf.Global.WxGzhConf.AppSecret
	token, err := GetWxGzhAccessToken(ctx, appId, appSecret)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	resp, err := restysdk.NewRequest().SetQueryParams(map[string]string{
		"access_token": token,
		"media_id":     mediaId,
	}).Get(constant.WxGzhMediaUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return resp.Body(), nil
}

// CheckWxGzhApiErr 检查是否错误
func CheckWxGzhApiErr(data []byte) error {
	var apiErr constant.WxGzhApiErr
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
