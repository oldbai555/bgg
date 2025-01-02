/**
 * @Author: zjj
 * @Date: 2025/1/7
 * @Desc:
**/

package deepseek

import (
	"github.com/oldbai555/bgg/pkg/syscfg"
	"testing"
)

func TestChatCompletions(t *testing.T) {
	seek, err := syscfg.GetDeepSeek()
	if err != nil {
		t.Logf("err:%v", err)
		return
	}
	completions, err := ChatCompletions(&ChatChatCompletionsReq{
		BaseUrl: seek.BaseUrl,
		Token:   seek.Token,
		MsgList: []*ChatCompletionsMessage{
			{
				Content: "你好，我有一个问题想要问你",
				Role:    "user",
			},
		},
	})
	if err != nil {
		t.Logf("err:%v", err)
		return
	}
	for _, message := range completions.MsgList {
		t.Logf("message:%s\n", message.Content)
	}
}
