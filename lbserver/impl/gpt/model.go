package gpt

import "fmt"

type ChatCompletionReq struct {
	Model    string                      `json:"model"`
	Messages []*ChatCompletionReqMessage `json:"messages"`
}

type ChatCompletionReqMessage struct {
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

// APIError represents an error that occurred on an API
type APIError struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"`
	Code    interface{} `json:"code"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%v:%s] %s , param:%v", e.Code, e.Type, e.Message, e.Param)
}

// APIErrorResponse is the full error response that has been returned by an API.
type APIErrorResponse struct {
	Error APIError `json:"error"`
}
