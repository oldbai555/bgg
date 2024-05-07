/**
 * @Author: zjj
 * @Date: 2024/5/7
 * @Desc:
**/

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/service/lboss/compress"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/json"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bgin"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

func registerRouter(r *gin.Engine) {
	r.MaxMultipartMemory = MaxMultipartMemory

	//r.StaticFS("/files", http.Dir(BaseStoragePath))
	r.StaticFS("/js", http.Dir(BaseJsPath))
	r.LoadHTMLGlob(BaseTemplatesPath)
	r.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/view/index"
		r.HandleContext(ctx)
	})

	group := r.Group("view")
	group.GET("index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	group = r.Group("/api")
	group.POST("/upload", handleUpload)
	group.GET("/download/*sUrl", handleDownload)
	group.GET("/syncfileindex", handleSyncFileIndex)
	group.GET("/sortUrlList", handleSortUrlList)
	group.GET("/clean", handleClean)
}

func handleUpload(ctx *gin.Context) {
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
	filePath := path.Join(BaseStoragePath, timeFmt)
	savePath := path.Join(filePath, fileName)

	// 判断如果是windows环境 则需要将 savePath FromSlash
	savePath = ToSlash(savePath)

	// 3.判断文件是否存在
	if _, err = os.Stat(savePath); err == nil {
		handler.Error(lberr.NewErr(500, "保存文件失败，文件重复"))
		return
	}

	// 4.创建存储路径文件夹
	if !utils.FileExists(filePath) {
		err = os.MkdirAll(filePath, 0775)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(lberr.NewErr(500, "创建文件夹失败"))
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
	var fileInfo File
	fileInfo.Name = file.Filename
	fileInfo.ReName = fileName
	fileInfo.Path = savePath
	fileInfo.Md5 = GetFileMd5(saveFile)
	fileInfo.Size = file.Size
	fileInfo.TimeStamp = time.Now().UnixNano() / 1e6

	// 7.保存文件路径和索引到数据库
	fileInfoJson, err := json.Marshal(fileInfo)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(lberr.NewErr(500, "保存文件失败，写入文件异常"))
		return
	}

	sUrl := compress.GenShortUrl(compress.CharsetRandomAlphanumeric, savePath, func(url, keyword string) bool {
		data, _ := dbConn.Get([]byte(keyword), nil)
		if data == nil {
			return true
		}
		return false
	})
	err = dbConn.Put([]byte(sUrl), fileInfoJson, nil)

	// 7.返回文件唯一索引
	handler.HttpJson(sUrl)
}

func handleDownload(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	param := ctx.Param("sUrl")
	sUrl := strings.TrimPrefix(param, "/")
	if sUrl == "" {
		handler.Error(lberr.NewErr(500, "下载文件失败，参数错误"))
		return
	}

	data, err := dbConn.Get([]byte(sUrl), nil)
	if err != nil {
		handler.Error(err)
		return
	}

	var fileInfo File
	if err = json.Unmarshal(data, &fileInfo); err != nil {
		handler.Error(err)
		return
	}
	fileName := url.QueryEscape(fileInfo.Name)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.File(fileInfo.Path)
}

// 同步文件索引
func handleSyncFileIndex(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)

	sUrlList, err := syncFileIndex()
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	handler.HttpJson(map[string]interface{}{
		"sortKeys": sUrlList,
	})
}

func handleSortUrlList(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	iter := dbConn.NewIterator(nil, nil)
	var retMap = make(map[string]*File)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		var newFile File
		err := json.Unmarshal(value, &newFile)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}
		retMap[string(key)] = &newFile
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	handler.HttpJson(retMap)
}

func handleClean(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	iter := dbConn.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		err := dbConn.Delete(key, nil)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	handler.HttpJson("")
}
