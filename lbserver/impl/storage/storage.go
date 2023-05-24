package storage

import (
	"bytes"
	webtool "github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/storage"
	"net/http"
)

var S storage.FileStorageInterface

func InitStorage(conf *webtool.StorageConf) {
	storage.Setup(storage.Config{
		Type:      conf.Type,
		SecretID:  conf.SecretId,
		SecretKey: conf.SecretKey,
		BucketURL: conf.BucketUrl,
	})
	S = storage.FileStorage
}

// ConvertMediaUrl 将URL转换为系统的URL
func ConvertMediaUrl(filename, url string) (string, error) {
	objectKey := `public/link-info/assets/images/` + filename

	open, err := http.Get(url)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	err = S.Put(objectKey, open.Body)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	signedURL, err := S.SignURL(objectKey, http.MethodGet, 60*60*24*365)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	return signedURL, nil
}

// ConvertMediaBytes 将字节流给上传成URL
func ConvertMediaBytes(filename string, b []byte) (string, error) {
	objectKey := `public/link-info/assets/images/` + filename
	buffer := bytes.NewBuffer(b)
	err := S.Put(objectKey, buffer)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	signedURL, err := S.SignURL(objectKey, http.MethodGet, 60*60*24*365)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	return signedURL, nil
}
