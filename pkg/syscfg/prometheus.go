package syscfg

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const defaultPrometheusPrefix = "prometheus"

type PrometheusConf struct {
	Ip   string `json:"ip"`
	Port uint32 `json:"port"`
}

func NewPrometheusConf(viper *viper.Viper) *PrometheusConf {
	var v PrometheusConf
	val := viper.Get(defaultPrometheusPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}
