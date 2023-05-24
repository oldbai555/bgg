package wordscheck

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/oldbai555/lbtool/log"
	"strings"
)

type WordData struct {
	Keyword  string `json:"keyword"`
	Category string `json:"category"`
	Position string `json:"position"`
}

type CheckResp struct {
	Code uint32     `json:"code"`
	Msg  string     `json:"msg"`
	Data []WordData `json:"data"`
}

// DoCheckWords true 有敏感内容
func DoCheckWords(url, content string) (bool, error) {

	// 敏感词枚举
	if strings.Contains(content, "色情") ||
		strings.Contains(content, "暴力") ||
		strings.Contains(content, "党") ||
		strings.Contains(content, "党中央") ||
		strings.Contains(content, "中国") ||
		strings.Contains(content, "台独") ||
		strings.Contains(content, "港独") ||
		strings.Contains(content, "杀人") ||
		strings.Contains(content, "抢劫") ||
		strings.Contains(content, "强奸") ||
		strings.Contains(content, "赌博") ||
		strings.Contains(content, "彩票") ||
		strings.Contains(content, "六合彩") ||
		strings.Contains(content, "黄色网站") {
		return true, nil
	}

	values := map[string]interface{}{
		"isBiz":   1,
		"content": content,
	}

	resp, err := resty.New().NewRequest().SetBody(values).Post(url)
	if err != nil {
		log.Errorf("err:%v", err)
		return false, err
	}

	var checkResp CheckResp
	err = json.Unmarshal(resp.Body(), &checkResp)
	if err != nil {
		log.Errorf("err:%v", err)
		return false, err
	}

	log.Infof("code=%d", checkResp.Code)
	log.Infof("msg=%s", checkResp.Msg)
	log.Infof("data=%+v", checkResp.Data)

	return len(checkResp.Data) > 0, nil
}
