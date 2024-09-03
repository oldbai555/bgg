package syscfg

import (
	"encoding/json"
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

type SysCfg struct {
	V *viper.Viper

	StorageConf *StorageConf
	ServerConf  *ServerConf
	Proxys      *Proxys
}

var sc *SysCfg

func New(viper *viper.Viper, option ...Option) (*SysCfg, error) {
	sc = &SysCfg{
		V: viper,
	}
	option = append(option, OptionWithServer())
	// 初始化组件
	for _, o := range option {
		o(sc)
	}
	return sc, nil
}

func JsonConvertStruct(re interface{}, out interface{}) error {
	marshal, err := json.Marshal(re)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}

	err = json.Unmarshal(marshal, out)
	if err != nil {
		log.Errorf("err is : %v", err)
		return err
	}
	return nil
}
