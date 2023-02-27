package impl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/oldbai555/bgg/client/lbuser"
	webtool2 "github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/restysdk"
	gogpt "github.com/sashabaranov/go-gpt3"
	"strings"
	"testing"
)

func init() {
	v, err := initViper()
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}
	lb = &Tool{}
	lb.WebTool, _ = webtool2.NewWebTool(v, webtool2.OptionWithOrm(&lbuser.ModelUser{}), webtool2.OptionWithRdb())
	InitDbOrm()
}

var dataPrefix = []byte("data: ")
var doneSequence = []byte("[DONE]")

func TestLbwebsocketServer_HandleWs(t *testing.T) {
	const apiURLv1 = "https://api.openai.com/v1"
	urlSuffix := "/completions"

	// 接入chatGPT
	resp, err := restysdk.NewRequest().SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "text/event-stream",
		"Cache-Control": "no-cache",
		"Connection":    "keep-alive",
		"Authorization": fmt.Sprintf("Bearer %s", lb.V.GetString("chatGpt.api_key")),
	}).SetBody(gogpt.CompletionRequest{
		Prompt:    "如何快速赚钱",
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 3000,
		Stream:    true,
	}).Post(fmt.Sprintf("%s%s", apiURLv1, urlSuffix))
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

	err = checkForSuccess(resp)
	if err != nil {
		log.Errorf("err is %v", err)
		return
	}

	var respText string
	reader := bufio.NewReader(bytes.NewReader(resp.Body()))
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Errorf("err is %v", err)
			return
		}

		// make sure there isn't any extra whitespace before or after
		line = bytes.TrimSpace(line)

		// the completion API only returns data events
		if !bytes.HasPrefix(line, dataPrefix) {
			continue
		}
		line = bytes.TrimPrefix(line, dataPrefix)

		// the stream is completed when terminated by [DONE]
		if bytes.HasPrefix(line, doneSequence) {
			break
		}
		output := new(gogpt.CompletionResponse)
		if err := json.Unmarshal(line, output); err != nil {
			log.Errorf("invalid json stream data: %v", err)
			return
		}
		respText = respText + output.Choices[0].Text
	}
	respText = strings.TrimSpace(strings.TrimPrefix(respText, "?"))
	log.Infof("resp is %s", respText)
}

// returns an error if this response includes an error.
func checkForSuccess(resp *resty.Response) error {
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}

	var result APIErrorResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := APIError{
			Code:    resp.StatusCode(),
			Type:    "Unexpected",
			Message: string(resp.Body()),
		}
		return apiError
	}
	result.Error.Code = resp.StatusCode()
	return result.Error
}

type APIErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param"`
	Code    interface{} `json:"code"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%d:%s] %s", e.Code, e.Type, e.Message)
}
