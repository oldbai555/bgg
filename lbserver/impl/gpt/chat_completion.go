package gpt

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/oldbai555/lbtool/log"
	"time"
)

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
	}).SetBody(req).Post(fmt.Sprintf("%s%s", ApiURLv1, urlSuffix))
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
