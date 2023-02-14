package impl

import (
	"context"
	"errors"
	"github.com/oldbai555/bgg/lbwebsocket"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"time"
)

var lbwebsocketServer LbwebsocketServer

type LbwebsocketServer struct {
	*lbwebsocket.UnimplementedLbwebsocketServer
}

func (a *LbwebsocketServer) HandleWs(ctx context.Context, req *lbwebsocket.HandleWsReq) (*lbwebsocket.HandleWsRsp, error) {
	var rsp lbwebsocket.HandleWsRsp

	// 先去检查会话是否存在
	var chat lbwebsocket.ModelChat
	err := ChatOrm.NewScope().UpdateOrCreate(ctx, map[string]interface{}{
		lbwebsocket.FieldVid_:      req.WsData.Vid,
		lbwebsocket.FieldUsername_: "oldbai",
	}, map[string]interface{}{
		lbwebsocket.FieldVid_:      req.WsData.Vid,
		lbwebsocket.FieldUsername_: "oldbai",
	}, &chat)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	// 访客信息
	var visitor lbwebsocket.ModelVisitor
	err = VisitorOrm.NewScope().UpdateOrCreate(ctx,
		map[string]interface{}{
			lbwebsocket.FieldVid_: req.WsData.Vid,
		},
		map[string]interface{}{
			lbwebsocket.FieldName_: req.WsData.VName,
			lbwebsocket.FieldId_:   req.Ip,
		},
		&visitor)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	// 消息
	err = MessageOrm.NewScope().Create(ctx, &lbwebsocket.ModelMessage{
		Content:      req.WsData.Content,
		Cid:          chat.Id,
		From:         visitor.Vid,
		Receive:      "oldbai",
		IsFromSys:    false,
		SysMessageId: req.WsData.MsgId,
		SendAt:       req.WsData.SendAt,
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	// 接入chatGPT
	var gpt = gogpt.NewClient("")
	stream, err := gpt.CreateCompletionStream(ctx, gogpt.CompletionRequest{
		Prompt:    req.WsData.Content.Text.Text,
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 2048,
		Stream:    true,
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	var respText string
	for {
		resp, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Errorf("err is %v", err)
			break
		}
		if len(resp.Choices) == 0 {
			continue
		}
		respText = respText + resp.Choices[0].Text
	}

	// 封装回复的消息
	newContent := &lbwebsocket.Content{
		Type: uint32(lbwebsocket.Content_TypeText),
		Text: &lbwebsocket.Content_Text{
			Text: respText,
		},
	}
	uuid := utils.GenUUID()
	err = MessageOrm.NewScope().Create(ctx, &lbwebsocket.ModelMessage{
		Content:      newContent,
		Cid:          chat.Id,
		From:         "oldbai",
		Receive:      visitor.Vid,
		IsFromSys:    true,
		SysMessageId: uuid,
		SendAt:       uint64(time.Now().UnixMilli()),
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	rsp.WsData = &lbwebsocket.WsData{
		Content:  newContent,
		Vid:      visitor.Vid,
		Username: "oldbai",
		Type:     uint32(lbwebsocket.WsType_WsTypeMessage),
		VName:    req.WsData.VName,
		UName:    req.WsData.UName,
	}
	return &rsp, nil
}
