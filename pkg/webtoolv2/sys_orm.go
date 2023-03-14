package webtool

import (
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const defaultApolloMysqlPrefix = "mysql"
const defaultDatabase = "biz"

type GormMysqlConf struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewGormMysqlConf(viper *viper.Viper) *GormMysqlConf {
	var v GormMysqlConf
	val := viper.Get(defaultApolloMysqlPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}

func (m *GormMysqlConf) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.Username, m.Password, m.Addr, m.Port, defaultDatabase)
}
