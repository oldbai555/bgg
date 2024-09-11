package syscfg

import (
	"github.com/oldbai555/lbtool/pkg/json"
	"os"
	path2 "path"
	"path/filepath"
)

const defaultApolloRedisPrefix = "redis"

type RedisConf struct {
	Database int    `json:"database"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func NewRedisConf(path string) *RedisConf {
	if path == "" {
		path = defaultRedisConfPath
	}
	var v RedisConf
	data, err := os.ReadFile(filepath.ToSlash(path2.Join(path, "redis.json")))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &v)
	if err != nil {
		panic(err)
	}
	return &v
}
