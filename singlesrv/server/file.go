/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package server

import (
	"context"
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/compress"
	"github.com/oldbai555/bgg/pkg/tool"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/bgg/singlesrv/server/constant"
	"github.com/oldbai555/bgg/singlesrv/server/mq"
	"github.com/oldbai555/bgg/singlesrv/server/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/micro/uctx"
	"os"
)

func SyncFileIndex() error {
	err := tool.ListFile(constant.BaseStoragePath, func(path string, info os.FileInfo) {
		file, _ := os.Open(path)
		path = tool.ToSlash(path)
		var exist bool
		sUrl := compress.GenShortUrl(compress.CharsetRandomAlphanumeric, path, func(url, keyword string) bool {
			if exist {
				return false
			}
			_, err := cache.GetFileBySortUrl(keyword)
			if err != nil && !cache.IsNotFound(err) {
				log.Errorf("err:%v", err)
				return true
			}
			if cache.IsNotFound(err) {
				return true
			}
			exist = true
			return false
		})
		if exist {
			return
		}
		if sUrl == "" {
			return
		}
		pubErr := MqTopicBySyncFile.Pub(uctx.NewBaseUCtx(), &client.ModelFile{
			Size:    info.Size(),
			Name:    info.Name(),
			Rename:  info.Name(),
			Path:    path,
			Md5:     tool.GetFileMd5(file),
			SortUrl: sUrl,
		})
		if pubErr != nil {
			log.Errorf("err:%v", pubErr)
		}
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func MqTopicBySyncFileHandler(msg *nsq.Message) error {
	return mq.Process(msg, func(buf []byte) error {
		var data client.ModelFile
		err := MqTopicBySyncFile.Unmarshal(buf, &data)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if data.Md5 == "" {
			return client.ErrFileMd5IsEmpty
		}
		db := mysql.File.NewScope(context.Background())
		err = db.Eq(client.FieldMd5_, data.Md5).First(&client.ModelFile{})
		if err != nil && !mysql.File.IsNotFoundErr(err) {
			log.Errorf("err:%v", err)
			return err
		}
		if err == nil {
			return client.ErrFileMd5Already
		}
		_, err = mysql.File.NewScope(context.Background()).Create(&data)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}
