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

// SyncFileIndex 同步最新文件索引 - 先临时冗余设计 后续在解耦优化
func SyncFileIndex(ctx context.Context, cacheAllFile bool) error {
	var filePathMap = make(map[string]string)
	var sortUrlMap = make(map[string]struct{})

	var list []*client.ModelFile
	err := mysql.File.NewScope(ctx).Chunk(2000, &list, func() (stop bool) {
		for _, file := range list {
			filePathMap[file.Path] = file.Md5
			sortUrlMap[file.SortUrl] = struct{}{}
		}
		return
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var newFileList []*client.ModelFile

	// 同步本地文件
	err = tool.ListFile(constant.BaseStoragePath, func(path string, info os.FileInfo) {
		file, err := os.Open(path)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		defer func() {
			if file == nil {
				return
			}
			err := file.Close()
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		}()
		path = tool.ToSlash(path)

		// 校验文件是否在文件夹中
		_, ok := filePathMap[path]
		if ok {
			delete(filePathMap, path)
			return
		}

		var canUse bool
		sUrl := compress.GenShortUrl(compress.CharsetRandomAlphanumeric, path, func(url, keyword string) bool {
			_, ok := sortUrlMap[keyword]
			if ok {
				canUse = false
				return canUse
			}
			canUse = true
			return canUse
		})
		if !canUse {
			return
		}
		if sUrl == "" {
			return
		}
		newFileList = append(newFileList, &client.ModelFile{
			Size:    info.Size(),
			Name:    info.Name(),
			Rename:  info.Name(),
			Path:    path,
			Md5:     tool.GetFileMd5(file),
			SortUrl: sUrl,
		})
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		uCtx = uctx.NewBaseUCtx()
	}
	if len(newFileList) != 0 {
		pubErr := MqTopicSyncFile.Pub(uCtx, &client.MqSyncFile{
			FileList: newFileList,
		})
		if pubErr != nil {
			log.Errorf("err:%v", pubErr)
		}
	}

	// 打上标记
	if len(filePathMap) != 0 {
		var md5List []string
		for _, md5 := range filePathMap {
			md5List = append(md5List, md5)
		}
		_, err = mysql.File.NewScope(ctx).
			Eq(client.FieldState_, uint32(client.ModelFile_StateNil)).
			In(client.FieldMd5_, md5List).
			Update(map[string]interface{}{
				client.FieldState_: uint32(client.ModelFile_StateNotFoundInLocalDir),
			})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		_, err = mysql.File.NewScope(ctx).
			Eq(client.FieldState_, uint32(client.ModelFile_StateNotFoundInLocalDir)).
			NotIn(client.FieldMd5_, md5List).
			Update(map[string]interface{}{
				client.FieldState_: uint32(client.ModelFile_StateNil),
			})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	if cacheAllFile {
		// 避免报错 传个空值
		err = MqTopicCacheAllFile.DeferredPublish(uCtx, time.Second*3, &client.MqCacheAllFile{})
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func MqTopicBySyncFileHandler(msg *nsq.Message) error {
	var doSingleFileLogic = func(file *client.ModelFile) error {
		if file == nil {
			log.Errorf("unmarshal file is nil")
			return nil
		}

		if file.Md5 == "" {
			return client.ErrFileMd5IsEmpty
		}

		db := mysql.File.NewScope(context.Background())
		err := db.Eq(client.FieldMd5_, file.Md5).First(&client.ModelFile{})
		if err != nil && !mysql.File.IsNotFoundErr(err) {
			log.Errorf("err:%v", err)
			return err
		}

		// 重新存一下缓存
		if err == nil {
			buf, err := MqTopicSyncFile.Marshal(file)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			return cache.SetFileBySortUrl(file.SortUrl, string(buf))
		}

		_, err = mysql.File.NewScope(context.Background()).Create(&file)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		buf, err := MqTopicSyncFile.Marshal(file)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = cache.SetFileBySortUrl(file.SortUrl, string(buf))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return nil
	}

	return mq.Process(msg, func(buf []byte) error {
		var data client.MqSyncFile
		err := MqTopicSyncFile.Unmarshal(buf, &data)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		for _, file := range data.FileList {
			err = doSingleFileLogic(file)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}

		return nil
	})
}

func MqTopicByCacheAllFileHandler(msg *nsq.Message) error {
	return mq.Process(msg, func(buf []byte) error {
		var list []*client.ModelFile
		err := mysql.File.NewScope(context.Background()).Chunk(2000, &list, func() (stop bool) {
			for _, file := range list {
				bytes, err := MqTopicCacheAllFile.Marshal(file)
				if err != nil {
					log.Errorf("err:%v", err)
				} else {
					err = cache.SetFileBySortUrl(file.SortUrl, string(bytes))
				}
			}
			return
		})
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
	err = MqTopicSyncFile.Pub(uctx.NewBaseUCtx(), fileInfo)
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
	if err = MqTopicSyncFile.Unmarshal([]byte(fileInfoStr), &fileInfo); err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	fileName := url.QueryEscape(fileInfo.Name)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.File(fileInfo.Path)
}
