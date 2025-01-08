/**
 * @Author: zjj
 * @Date: 2025/1/7
 * @Desc:
**/

package deepseek

import (
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	url2 "net/url"
)

// https://api-docs.deepseek.com/zh-cn/api/create-chat-completion

type ChatChatCompletionsReq struct {
	BaseUrl string
	Token   string
	MsgList []*ChatCompletionsMessage
}

type ChatCompletionsResp struct {
	MsgList []*ChatCompletionsMessage
}

func ChatCompletions(req *ChatChatCompletionsReq) (*ChatCompletionsResp, error) {
	// 先写死
	url, _ := url2.JoinPath(req.BaseUrl, "/chat/completions")

	var apiReq = &ChatCompletionsApiReq{
		Messages:         req.MsgList,
		Model:            "deepseek-chat",
		FrequencyPenalty: 0,
		MaxTokens:        2048,
		ResponseFormat: &ChatCompletionsResponseFormat{
			Type: "text",
		},
		Stop:          []string{"白哥哥"},
		StreamOptions: nil,
		Temperature:   1,
		TopP:          1,
		ToolChoice:    "auto",
	}

	request := restysdk.NewRequest()
	response, err := request.SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer " + req.Token,
	}).SetBody(apiReq).Post(url)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	var apiResp ChatCompletionsApiResp
	err = json.Unmarshal(response.Body(), &apiResp)
	if err != nil {
		return nil, lberr.Wrap(err)
	}

	var resp = &ChatCompletionsResp{
		MsgList: []*ChatCompletionsMessage{},
	}
	for _, choice := range apiResp.Choices {
		resp.MsgList = append(resp.MsgList, &ChatCompletionsMessage{
			Content: choice.Message.Content,
			Role:    choice.Message.Role,
		})
	}

	return resp, nil
}

type ChatCompletionsApiReq struct {
	Messages         []*ChatCompletionsMessage      `json:"messages"`
	Model            string                         `json:"model"`
	FrequencyPenalty int                            `json:"frequency_penalty"`
	MaxTokens        int                            `json:"max_tokens"`
	PresencePenalty  int                            `json:"presence_penalty"`
	ResponseFormat   *ChatCompletionsResponseFormat `json:"response_format"`
	Stop             []string                       `json:"stop"`
	Stream           bool                           `json:"stream"`
	StreamOptions    *ChatCompletionsStreamOptions  `json:"stream_options"`
	Temperature      int                            `json:"temperature"`
	TopP             int                            `json:"top_p"`
	ToolChoice       string                         `json:"tool_choice"`
}

type ChatCompletionsMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatCompletionsResponseFormat struct {
	Type string `json:"type"`
}

type ChatCompletionsStreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type ChatCompletionsApiResp struct {
	Id                string                   `json:"id"`
	Object            string                   `json:"object"`
	Created           int                      `json:"created"`
	Model             string                   `json:"model"`
	Choices           []*ChatCompletionsChoice `json:"choices"`
	UsAge             *ChatCompletionsUsAge    `json:"us age"`
	SystemFingerprint string                   `json:"system_fingerprint"`
}

type ChatCompletionsChoice struct {
	Index        int                     `json:"index"`
	Message      *ChatCompletionsMessage `json:"message"`
	FinishReason string                  `json:"finish_reason"`
}

type ChatCompletionsUsAge struct {
	PromptTokens          int `json:"prompt_tokens"`
	CompletionTokens      int `json:"completion_tokens"`
	TotalTokens           int `json:"total_tokens"`
	PromptCacheHitTokens  int `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int `json:"prompt_cache_miss_tokens"`
}
