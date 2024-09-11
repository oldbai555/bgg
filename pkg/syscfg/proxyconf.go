/**
 * @Author: zjj
 * @Date: 2024/9/3
 * @Desc:
**/

package syscfg

import "github.com/spf13/viper"

type Proxys struct {
	Map map[string]string `yaml:"map"`
}

const defaultProxysPrefix = "proxys"

func NewProxys(viper *viper.Viper) *Proxys {
	var v Proxys
	val := viper.Get(defaultProxysPrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		panic(err)
	}
	return &v
}
