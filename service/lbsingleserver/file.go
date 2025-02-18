/**
 * @Author: zjj
 * @Date: 2024/6/13
 * @Desc:
**/

package lbsingleserver

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"github.com/oldbai555/bgg/pkg/bctx"
	"github.com/oldbai555/bgg/pkg/constant"
	"github.com/oldbai555/bgg/pkg/tool"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/mq"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/lberr"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bgin"
	"github.com/oldbai555/micro/uctx"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// MqTopicBySyncFileHandler 消息队列-保存文件
func MqTopicBySyncFileHandler(msg *nsq.Message) error {
	return mq.Process[*lbsingle.MqSyncFile](msg, func(ctx uctx.IUCtx, data *lbsingle.MqSyncFile) error {
		for _, file := range data.FileList {
			err := saveFileToOrm(ctx, file)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}
		return nil
	})
}

// 预签名
func handlePreSigned(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	if minIoSDK == nil {
		return
	}
	nCtx := bctx.NewCtx(ctx, bctx.WithGinHeaderAuthorization(ctx))
	_, err := CheckAuth(nCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	object, err := minIoSDK.PreSignedPutObject(constant.BucketByPublic, ctx.Param("fileName"))
	if err != nil {
		handler.Error(err)
		return
	}
	handler.HttpJson(object)
}

// 上传文件
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

	var fileType uint32
	if isImageByExtension(file.Filename) {
		fileType = uint32(lbsingle.ModelFile_TypeImage)
	}

	sortUrl := utils.StrMd5(file.Filename)
	reFileName := sortUrl + path.Ext(file.Filename)
	p, err := minIoSDK.UploadNetIO(constant.BucketByPublic, reFileName, open)
	err = saveFileToOrm(nCtx, &lbsingle.ModelFile{
		Size:       file.Size,
		Name:       file.Filename,
		Rename:     reFileName,
		BucketPath: constant.BucketPath,
		Path:       p,
		SortUrl:    sortUrl,
		Md5:        sortUrl,
		Type:       fileType,
	})
	if err != nil {
		handler.Error(err)
		return
	}
	result, err := url.JoinPath(constant.BucketPath, sortUrl)
	if err != nil {
		handler.Error(err)
		return
	}
	handler.HttpJson(result)
}

// 上传包
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

// 下载文件
func handleDownloadFile(ctx *gin.Context) {
	handler := bgin.NewHandler(ctx)
	param := ctx.Param("sUrl")
	sUrl := strings.TrimPrefix(param, "/")
	if sUrl == "" {
		handler.Error(lberr.NewErr(500, "下载文件失败，参数错误"))
		return
	}
	nCtx := bctx.NewCtx(ctx)
	modelFile, err := OrmFile.NewBaseScope().Where(lbsingle.FieldSortUrl_, sUrl).First(nCtx)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(err)
		return
	}
	u, err := minIoSDK.PreSignedGetObject(constant.BucketByPublic, modelFile.Rename)
	if err != nil {
		log.Errorf("err:%v", err)
		handler.Error(lbsingle.ErrFileNotFound)
		return
	}
	ctx.Redirect(301, u)
}

// 保存到数据库
func saveFileToOrm(ctx uctx.IUCtx, file *lbsingle.ModelFile) error {
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

	err = OrmFile.NewBaseScope().Create(ctx, &file)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
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

// 保存文件至本地
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
	sUrl := utils.GenUUID()

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

	// 保存到数据库
	err = saveFileToOrm(ctx, fileInfo)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	return url.JoinPath(fileInfo.BucketPath, fileInfo.SortUrl)
}
