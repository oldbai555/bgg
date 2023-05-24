package webtool

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/spf13/viper"
)

const defaultApolloStoragePrefix = "storage"

type StorageConf struct {
	Type      string `json:"type"`
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	BucketUrl string `json:"bucket_url"`
}

func NewStorageConf(viper *viper.Viper) *StorageConf {
	var v StorageConf
	val := viper.Get(defaultApolloStoragePrefix)
	err := JsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		panic(err)
	}
	return &v
}
