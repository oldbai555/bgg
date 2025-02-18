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

type MinIOConf struct {
	Endpoint        string `json:"endpoint"`
	AccessKey       string `json:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func NewMinIOConf(path string) *MinIOConf {
	if path == "" {
		path = defaultConfMinIO
	}
	var v MinIOConf
	data, err := os.ReadFile(filepath.ToSlash(path2.Join(path, "minio.json")))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}
	return &v
}
