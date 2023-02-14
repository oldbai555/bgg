package webtool

import (
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/storage"
	"github.com/spf13/viper"
)

const defaultApolloStoragePrefix = "storage"

type StorageConf struct {
	Type      string `json:"type"`
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	BucketUrl string `json:"bucket_url"`
}

func (r *StorageConf) InitConf(viper *viper.Viper) error {
	var v StorageConf
	val := viper.Get(defaultApolloStoragePrefix)
	err := jsonConvertStruct(val, &v)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}
	log.Infof("init redis successfully")
	r.Type = v.Type
	r.SecretId = v.SecretId
	r.SecretKey = v.SecretKey
	r.BucketUrl = v.BucketUrl
	return nil
}

func (r *StorageConf) GenConfTool(tool *WebTool) error {
	log.Infof("init rdb engine successfully")
	storage.Setup(storage.Config{
		Type:      r.Type,
		SecretID:  r.SecretId,
		SecretKey: r.SecretKey,
		BucketURL: r.BucketUrl,
	})
	tool.Storage = storage.FileStorage
	return nil
}
