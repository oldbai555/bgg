/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/compress"
	"github.com/oldbai555/bgg/pkg/tool"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/cache"
	"github.com/oldbai555/bgg/singlesrv/server/constant"
	"github.com/oldbai555/bgg/singlesrv/server/mq"
	"github.com/oldbai555/bgg/singlesrv/server/mysql"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/uctx"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
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
			// 重新存一下缓存
			return cache.SetFileBySortUrl(data.SortUrl, string(buf))
		}
		_, err = mysql.File.NewScope(context.Background()).Create(&data)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		err = cache.SetFileBySortUrl(data.SortUrl, string(buf))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	})
}

func handleUploadFile(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)

	// 1.获取文件信息
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 2.构造文件存储路径, 可以很方便的按照天进行数据同步
	fileName := utils.Md5(utils.GenUUID()) + path.Ext(file.Filename)
	timeFmt := time.Unix(time.Now().Unix(), 0).Format("20060102")
	filePath := path.Join(constant.BaseStoragePath, timeFmt)
	savePath := path.Join(filePath, fileName)

	// 判断如果是windows环境 则需要将 savePath FromSlash
	savePath = tool.ToSlash(savePath)

	// 3.判断文件是否存在
	if _, err = os.Stat(savePath); err == nil {
		handler.Error(client.ErrFileAlreadyExist)
		return
	}

	// 4.创建存储路径文件夹
	if !utils.FileExists(filePath) {
		err = os.MkdirAll(filePath, 0775)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}
	}

	// 5.保存文件到本地目录
	err = ctx.SaveUploadedFile(file, savePath)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	saveFile, err := os.Open(savePath)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 6.构造保存结构
	var fileInfo = &client.ModelFile{
		Size:    file.Size,
		Name:    file.Filename,
		Rename:  fileName,
		Path:    savePath,
		Md5:     tool.GetFileMd5(saveFile),
		SortUrl: "",
	}

	sUrl := compress.GenShortUrl(compress.CharsetRandomAlphanumeric, savePath, func(url, keyword string) bool {
		_, err := cache.GetFileBySortUrl(keyword)
		if err != nil && !cache.IsNotFound(err) {
			log.Errorf("err:%v", err)
			return true
		}
		if cache.IsNotFound(err) {
			return true
		}
		return false
	})
	if sUrl == "" {
		handler.Error(client.ErrFileUploadFailure)
		return
	}
	fileInfo.SortUrl = sUrl

	// 7.交给MQ去保存
	err = MqTopicBySyncFile.Pub(uctx.NewBaseUCtx(), fileInfo)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 8.返回文件唯一索引
	handler.HttpJson(sUrl)
}

func handleDownloadFile(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	param := ctx.Param("sUrl")
	sUrl := strings.TrimPrefix(param, "/")
	if sUrl == "" {
		handler.Error(lberr.NewErr(500, "下载文件失败，参数错误"))
		return
	}

	fileInfoStr, err := cache.GetFileBySortUrl(sUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	var fileInfo client.ModelFile
	if err = MqTopicBySyncFile.Unmarshal([]byte(fileInfoStr), &fileInfo); err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	fileName := url.QueryEscape(fileInfo.Name)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.File(fileInfo.Path)
}
