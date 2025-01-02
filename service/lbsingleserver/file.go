/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package lbsingleserver

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/pkg/compress"
	"github.com/oldbai555/bgg/pkg/constant"
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
	"io"
	"net/http"
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

		// 新文件
		nCtx := bctx.NewCtx(context.Background())
		_, err = saveFile(nCtx, 0, constant.BaseStoragePath, info.Name(), file)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
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
		_, err = OrmFile.NewBaseScope().Update(ctx, map[string]interface{}{
			lbsingle.FieldSortUrl_: file.SortUrl,
		})
		if err != nil {
			log.Errorf("err:%v", err)
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
	nCtx := bctx.NewCtx(c, bctx.WithGinHeaderAuthorization(c))

	_, err := CheckAuth(nCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 获取文件信息
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	open, err := file.Open()
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	surl, err := saveFile(nCtx, 0, constant.BaseStoragePath, file.Filename, open)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 返回文件 URl
	handler.HttpJson(surl)
}

func handleUploadDeployFile(c *gin.Context) {
	handler := bgin.NewHandler(c)
	nCtx := bctx.NewCtx(c)

	// 获取文件信息
	file, err := c.FormFile("file")
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	open, err := file.Open()
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	surl, err := saveFile(nCtx, uint32(lbsingle.ModelFile_TypeSrvPack), constant.BaseDeployPath, file.Filename, open)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}

	// 返回文件 URl
	handler.HttpJson(surl)
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

func handleSyncFileIndex(c *gin.Context) {
	handler := bgin.NewHandler(c)
	err := SyncFileIndex(c, true)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	handler.HttpJson("ok")
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

// downloadFile 下载文件并保存到本地
func downloadFile(url string, basePath, fileName string) (string, error) {
	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Errorf("closeErr:%v", closeErr)
		}
	}()
	return saveFile(bctx.NewCtx(context.Background()), 0, basePath, fileName, resp.Body)
}

// 保存文件
func saveFile(ctx uctx.IUCtx, fileType uint32, basePath, fileName string, fileReader io.Reader) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, fileReader)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	size := buf.Len()

	if size == 0 {
		log.Errorf("file size is zero")
		return "", lbsingle.ErrFileNotFound
	}

	if fileType == 0 {
		if isImageByExtension(fileName) || isImageByMagicNumber(buf.Bytes()) {
			fileType = uint32(lbsingle.ModelFile_TypeImage)
		}
	}

	if basePath == "" {
		basePath = constant.BaseStoragePath
	}

	if fileName == "" {
		fileName = utils.GenUUID() + ".bin"
	}

	md5h := md5.New()
	md5h.Write(buf.Bytes())
	md5Str := fmt.Sprintf("%x", md5h.Sum(nil))

	var filePath, reFileName string
	switch fileType {
	case uint32(lbsingle.ModelFile_TypeSrvPack):
		filePath = basePath
		reFileName = fileName
	default:
		timeFmt := time.Unix(time.Now().Unix(), 0).Format("20060102")
		filePath = path.Join(basePath, timeFmt)
		reFileName = md5Str + path.Ext(fileName)
	}

	// 判断如果是windows环境 则需要将 savePath ToSlash
	savePath := path.Join(filePath, reFileName)
	savePath = tool.ToSlash(savePath)

	file, err := OrmFile.NewBaseScope().Where(lbsingle.FieldMd5_, md5Str).First(ctx)
	if err != nil && !OrmFile.IsNotFoundErr(err) {
		log.Errorf("err:%v", err)
		return "", err
	}

	// 本地文件
	_, statErr := os.Stat(savePath)

	// 本地存在且也入库
	if statErr == nil && file != nil {
		return url.JoinPath(file.BucketPath, file.SortUrl)
	}

	// 本地不存在 但是入库
	if statErr != nil && file == nil {
		_, err = OrmFile.NewBaseScope().Where(lbsingle.FieldMd5_, md5Str).Delete(ctx)
		if err != nil && !OrmFile.IsNotFoundErr(err) {
			log.Errorf("err:%v", err)
			return "", err
		}
	}

	// 创建存储路径文件夹
	if !utils.FileExists(filePath) {
		err = os.MkdirAll(filePath, 0775)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}
	}

	// 生成短链
	result, _ := url.JoinPath(constant.BaseStoragePath, reFileName)
	sUrl := compress.GenShortUrl(compress.CharsetRandomAlphanumeric, result, func(url, keyword string) bool {
		_, err := cache.GetFileBySortUrl(keyword)
		switch {
		case err != nil && !cache.IsNotFound(err):
			log.Errorf("err:%v", err)
			_, err := OrmFile.NewBaseScope().Where(lbsingle.FieldSortUrl_, keyword).First(ctx)
			if OrmFile.IsNotFoundErr(err) {
				return true
			}
			return false
		case cache.IsNotFound(err):
			return true
		default:
			return false
		}
	})
	if sUrl == "" {
		log.Errorf("gen sort url failed , err:%v", lbsingle.ErrFileAlreadyExist)
		return "", lbsingle.ErrFileAlreadyExist
	}

	// 保存到本地
	out, err := os.Create(savePath)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	defer func() {
		closeErr := out.Close()
		if closeErr != nil {
			log.Errorf("closeErr:%v", closeErr)
		}
	}()

	// 将响应体复制到文件中
	_, err = io.Copy(out, &buf)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	// 6.构造保存结构
	var fileInfo = &lbsingle.ModelFile{
		Size:       int64(size),
		Name:       fileName,
		Rename:     reFileName,
		Path:       savePath,
		Md5:        md5Str,
		SortUrl:    sUrl,
		Type:       fileType,
		BucketPath: constant.BucketPath,
	}

	// 保存到数据库和缓存
	err = doSingleFileLogic(ctx, fileInfo)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	return url.JoinPath(fileInfo.BucketPath, fileInfo.SortUrl)
}
