package impl

import (
	"context"
	"errors"
	"github.com/oldbai555/bgg/lbuser"
	"github.com/oldbai555/bgg/webtool"
	"github.com/oldbai555/lbtool/log"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"testing"
)

func init() {
	lb = &Tool{}
	lb.WebTool, _ = webtool.NewWebTool(webtool.OptionWithOrm(&lbuser.ModelUser{}), webtool.OptionWithRdb())
	InitDbOrm()
}

func TestLbwebsocketServer_HandleWs(t *testing.T) {
	// 接入chatGPT
	var gpt = gogpt.NewClient(lb.V.GetString("chatGpt.api_key"))
	stream, err := gpt.CreateCompletionStream(context.Background(), gogpt.CompletionRequest{
		Prompt:    "hello",
		Model:     gogpt.GPT3TextDavinci002,
		MaxTokens: 2048,
		Stream:    true,
	})
	if err != nil {
		log.Errorf("err is %v", err)
		return
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
	log.Infof("resp is %v", respText)
}
