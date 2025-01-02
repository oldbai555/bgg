/**
 * @Author: zjj
 * @Date: 2025/1/7
 * @Desc:
**/

package syscfg

import (
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"os"
	path2 "path"
	"path/filepath"
)

type DeepSeek struct {
	Token   string `json:"api_key"`
	BaseUrl string `json:"base_url"`
}

func NewDeepSeek(path string) *DeepSeek {
	if path == "" {
		path = defaultConfPath
	}
	var v DeepSeek
	data, err := os.ReadFile(filepath.ToSlash(path2.Join(path, "deepseek.json")))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}
	return &v
}

func GetDeepSeek() (*DeepSeek, error) {
	if Global == nil {
		return nil, lberr.NewInvalidArg("not found Global")
	}
	if Global.DeepSeek == nil {
		return nil, lberr.NewInvalidArg("not found DeepSeek")
	}
	return Global.DeepSeek, nil
}
