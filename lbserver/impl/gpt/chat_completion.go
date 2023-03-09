package gpt

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/oldbai555/lbtool/log"
	"time"
)

type ChatCompletionReq struct {
	Model    string                      `json:"model"`
	Messages []*ChatcompletionreqMessage `json:"messages"`
}

type ChatcompletionreqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRsp struct {
	Id      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int                        `json:"created"`
	Choices []*ChatCompletionRspChoice `json:"choices"`
	Usage   *ChatCompletionRspUsage    `json:"usage"`
}

type ChatCompletionRspUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionRspChoice struct {
	Index        int                             `json:"index"`
	Message      *ChatCompletionRspChoiceMessage `json:"message"`
	FinishReason string                          `json:"finish_reason"`
}

type ChatCompletionRspChoiceMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// APIError represents an error that occured on an API
type APIError struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"`
	Code    interface{} `json:"code"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%v:%s] %s , param:%v", e.Code, e.Type, e.Message, e.Param)
}

// APIErrorResponse is the full error respnose that has been returned by an API.
type APIErrorResponse struct {
	Error APIError `json:"error"`
}

const apiURLv1 = "https://api.openai.com/v1"

// var defaultHeaders = map[string]string{
// 	"Content-type":  "application/json",
// 	"Accept":        "application/json; charset=utf-8",
// 	"Authorization": fmt.Sprintf("Bearer %s", apiKey),
// }

// returns an error if this response includes an error.
func checkForSuccess(data []byte, statusCode int) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	var result APIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := APIError{
			Code:    statusCode,
			Type:    "Unexpected",
			Message: string(data),
		}
		return apiError
	}

	result.Error.Code = statusCode
	return result.Error
}

func DoChatCompletion(proxy string, apiKey string, req *ChatCompletionReq) (*ChatCompletionRsp, error) {
	var rsp ChatCompletionRsp
	urlSuffix := "/chat/completions"
	client := resty.New()
	if proxy != "" {
		client.SetProxy(proxy)
	}
	resp, err := client.SetTimeout(time.Minute).NewRequest().SetHeaders(map[string]string{ // 最多等待一分钟的回复
		"Content-type":  "application/json",
		"Accept":        "application/json; charset=utf-8",
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
	}).SetBody(req).Post(fmt.Sprintf("%s%s", apiURLv1, urlSuffix))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = checkForSuccess(resp.Body(), resp.StatusCode())
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
