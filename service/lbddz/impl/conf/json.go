package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"os"
	"path/filepath"
)

var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
}

func init() {
	data, err := os.ReadFile(filepath.ToSlash("/Users/zhangjianjun/work/lb/github.com/oldbai555/bgg/service/lbddz/impl/conf/server.json"))
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
