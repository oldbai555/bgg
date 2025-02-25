/**
 * @Author: zjj
 * @Date: 2025/2/24
 * @Desc:
**/

package syscfg

import (
	"github.com/oldbai555/lbtool/pkg/json"
	"os"
	path2 "path"
	"path/filepath"
)

type NsqConf struct {
	Address string `json:"address"`
}

func NewNsqConf(path string) *NsqConf {
	if path == "" {
		path = defaultConfNsq
	}
	var v NsqConf
	data, err := os.ReadFile(filepath.ToSlash(path2.Join(path, "nsq.json")))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}
	return &v
}
