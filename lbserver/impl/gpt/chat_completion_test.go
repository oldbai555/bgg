package gpt

import (
	"github.com/oldbai555/lbtool/log"
	"testing"
)

func TestDoChatCompletion(t *testing.T) {
	completionRsp, err := DoChatCompletion("", "", &ChatCompletionReq{
		Model: "gpt-3.5-turbo",
		Messages: []*ChatCompletionReqMessage{
			{
				Content: "早上好",
				Role:    "user",
			},
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	for _, choice := range completionRsp.Choices {
		log.Infof("msg is %v", choice.Message.Content)
	}
}
