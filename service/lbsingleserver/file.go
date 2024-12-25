/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package lbsingleserver

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/pkg/compress"
	constant2 "github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/bgg/pkg/marshal"
	"github.com/oldbai555/bgg/pkg/tool"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/cache"
	"github.com/oldbai555/bgg/service/lbsingleserver/mq"
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

	uCtx, err := uctx.ToUCtx(ctx)
	if err != nil {
		uCtx = uctx.NewBaseUCtx()
		uCtx.SetTraceId(utils.GenUUID())
	}

	err = OrmFile.NewBaseScope().Chunk(uCtx, 2000, func(out []*lbsingle.ModelFile) error {
		for _, file := range out {
			filePathMap[file.Path] = file.Md5
			sortUrlMap[file.SortUrl] = struct{}{}
		}
		return nil
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var newFileList []*lbsingle.ModelFile

	// 同步本地文件
	err = tool.ListFile(constant2.BaseStoragePath, func(path string, info os.FileInfo) {
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
		newFileList = append(newFileList, &lbsingle.ModelFile{
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

	for _, file := range newFileList {
		err := doSingleFileLogic(uCtx, file)
		if err != nil {
			log.Errorf("save file failed %s", file.Name)
		}
	}

	// 打上标记
	if len(filePathMap) != 0 {
		var md5List []string
		for _, md5 := range filePathMap {
			md5List = append(md5List, md5)
		}
		_, err = OrmFile.NewBaseScope().
			Where(lbsingle.FieldState_, uint32(lbsingle.ModelFile_StateNil)).
			WhereIn(lbsingle.FieldMd5_, md5List).
			Update(uCtx, map[string]interface{}{
				lbsingle.FieldState_: uint32(lbsingle.ModelFile_StateNotFoundInLocalDir),
			})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		_, err = OrmFile.NewBaseScope().
			Where(lbsingle.FieldState_, uint32(lbsingle.ModelFile_StateNotFoundInLocalDir)).
			WhereNotIn(lbsingle.FieldMd5_, md5List).
			Update(uCtx, map[string]interface{}{
				lbsingle.FieldState_: uint32(lbsingle.ModelFile_StateNil),
			})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	if cacheAllFile {
		// 避免报错 传个空值
		err = MqTopicCacheAllFile.DeferredPublish(uCtx, time.Second*3, &lbsingle.MqCacheAllFile{})
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func doSingleFileLogic(ctx uctx.IUCtx, file *lbsingle.ModelFile) error {
	if file == nil {
		log.Errorf("unmarshal file is nil")
		return nil
	}

	if file.Md5 == "" {
		return lbsingle.ErrFileMd5IsEmpty
	}
	db := OrmFile.NewBaseScope()
	_, err := db.Where(lbsingle.FieldMd5_, file.Md5).First(ctx)
	if err != nil && !OrmFile.IsNotFoundErr(err) {
		log.Errorf("err:%v", err)
		return err
	}

	// 重新存一下缓存
	if err == nil {
		buf, err := marshal.PbMarshal(file)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		return cache.SetFileBySortUrl(file.SortUrl, string(buf))
	}

	err = OrmFile.NewBaseScope().Create(ctx, &file)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	buf, err := marshal.PbMarshal(file)
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

func MqTopicBySyncFileHandler(msg *nsq.Message) error {

	return mq.Process[*lbsingle.MqSyncFile](msg, func(ctx uctx.IUCtx, data *lbsingle.MqSyncFile) error {
		for _, file := range data.FileList {
			err := doSingleFileLogic(ctx, file)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}

		return nil
	})
}

func MqTopicByCacheAllFileHandler(msg *nsq.Message) error {
	return mq.Process[*lbsingle.MqCacheAllFile](msg, func(ctx uctx.IUCtx, _ *lbsingle.MqCacheAllFile) error {
		err := OrmFile.NewBaseScope().Chunk(ctx, 2000, func(out []*lbsingle.ModelFile) error {
			for _, file := range out {
				bytes, err := marshal.PbMarshal(file)
				if err != nil {
					log.Errorf("err:%v", err)
				} else {
					err = cache.SetFileBySortUrl(file.SortUrl, string(bytes))
				}
			}
			return nil
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		return nil
	})
}

func handleUploadFile(c *gin.Context) {
	handler := bgin.NewHandler(c)
	nCtx := bctx.NewCtx(
		c,
		bctx.WithGinHeaderAuthorization(c),
	)
	_, err := CheckAuth(nCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 1.获取文件信息
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 2.构造文件存储路径, 可以很方便的按照天进行数据同步
	fileName := utils.Md5(utils.GenUUID()) + path.Ext(file.Filename)
	timeFmt := time.Unix(time.Now().Unix(), 0).Format("20060102")
	filePath := path.Join(constant2.BaseStoragePath, timeFmt)
	savePath := path.Join(filePath, fileName)

	// 判断如果是windows环境 则需要将 savePath ToSlash
	savePath = tool.ToSlash(savePath)

	// 3.判断文件是否存在
	if _, err = os.Stat(savePath); err == nil {
		handler.Error(lbsingle.ErrFileAlreadyExist)
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
	err = c.SaveUploadedFile(file, savePath)
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

	var typ = uint32(0)
	var fileBytes []byte
	_, err = saveFile.Read(fileBytes)
	if err != nil {
		return
	}
	if isImageByExtension(fileName) || isImageByMagicNumber(fileBytes) {
		typ = uint32(lbsingle.ModelFile_TypeImage)
	}

	// 6.构造保存结构
	var fileInfo = &lbsingle.ModelFile{
		Size:    file.Size,
		Name:    file.Filename,
		Rename:  fileName,
		Path:    savePath,
		Md5:     tool.GetFileMd5(saveFile),
		SortUrl: "",
		Type:    typ,
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
		handler.Error(lbsingle.ErrFileUploadFailure)
		return
	}
	fileInfo.SortUrl = sUrl

	// 7.保存
	var aSync = c.GetHeader(constant2.HEADER_LBSINGLE_ASYNC)
	var fileList = []*lbsingle.ModelFile{fileInfo}
	if aSync != "" {
		err = MqTopicSyncFile.Pub(nCtx, &lbsingle.MqSyncFile{
			FileList: fileList,
		})
	} else {
		err = doSingleFileLogic(nCtx, fileInfo)
	}
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 8.返回文件唯一索引
	handler.HttpJson(sUrl)
}

func handleUploadDeployFile(c *gin.Context) {
	handler := bgin.NewHandler(c)
	nCtx := bctx.NewCtx(
		c,
		bctx.WithGinHeaderAuthorization(c),
	)

	// 1.获取文件信息
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 2.构造文件存储路径
	fileName := path.Base(file.Filename)
	filePath := path.Join(constant2.BaseDeployPath)
	savePath := path.Join(filePath, fileName)

	// 判断如果是windows环境 则需要将 savePath ToSlash
	savePath = tool.ToSlash(savePath)

	// 3.判断文件是否存在
	if _, err = os.Stat(savePath); err == nil {
		handler.Error(lbsingle.ErrFileAlreadyExist)
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
	err = c.SaveUploadedFile(file, savePath)
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
	var fileInfo = &lbsingle.ModelFile{
		Size:    file.Size,
		Name:    file.Filename,
		Rename:  fileName,
		Path:    savePath,
		Md5:     tool.GetFileMd5(saveFile),
		SortUrl: "",
		Type:    uint32(lbsingle.ModelFile_TypeSrvPack),
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
		handler.Error(lbsingle.ErrFileUploadFailure)
		return
	}
	fileInfo.SortUrl = sUrl

	// 7.保存
	var aSync = c.GetHeader(constant2.HEADER_LBSINGLE_ASYNC)
	var fileList = []*lbsingle.ModelFile{fileInfo}
	if aSync != "" {
		err = MqTopicSyncFile.Pub(nCtx, &lbsingle.MqSyncFile{
			FileList: fileList,
		})
	} else {
		err = doSingleFileLogic(nCtx, fileInfo)
	}
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

	nCtx := bctx.NewCtx(ctx)
	fileInfoStr, err := cache.GetFileBySortUrl(sUrl)
	if err != nil && !cache.IsNotFound(err) {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	var filePath string
	var fileName string
	switch {
	case err != nil && !cache.IsNotFound(err):
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	case cache.IsNotFound(err):
		modelFile, err := OrmFile.NewBaseScope().Where(lbsingle.FieldSortUrl_, sUrl).First(nCtx)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}
		p := modelFile.Path
		_, err = os.Stat(p)
		if err != nil {
			log.Errorf("err:%v", err)
			handler.Error(lbsingle.ErrFileNotFound)
			return
		}

		// 重新存一下缓存
		if err == nil {
			buf, err := marshal.PbMarshal(modelFile)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			err = cache.SetFileBySortUrl(modelFile.SortUrl, string(buf))
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}

		fileName = url.QueryEscape(modelFile.Name)
		filePath = modelFile.Path
	default:
		var fileInfo lbsingle.ModelFile
		if err = marshal.PbUnmarshal([]byte(fileInfoStr), &fileInfo); err != nil {
			log.Errorf("err:%v", err)
			handler.Error(err)
			return
		}
		fileName = url.QueryEscape(fileInfo.Name)
		filePath = fileInfo.Path
	}

	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	ctx.File(filePath)
}

// 通过文件头信息（魔术字节）判断
func isImageByMagicNumber(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// 常见图片格式的魔术字节
	magicNumbers := map[string][]byte{
		"JPEG": {0xFF, 0xD8, 0xFF},
		"PNG":  {0x89, 0x50, 0x4E, 0x47},
		"GIF":  {0x47, 0x49, 0x46, 0x38},
		"BMP":  {0x42, 0x4D},
		"TIFF": {0x49, 0x49, 0x2A, 0x00},
		"WEBP": {0x52, 0x49, 0x46, 0x46},
	}

	for _, magic := range magicNumbers {
		if bytes.HasPrefix(data, magic) {
			return true
		}
	}
	return false
}

// 通过文件扩展名判断
func isImageByExtension(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".webp": true,
		// 添加其他图片格式
	}
	return imageExtensions[ext]
}
