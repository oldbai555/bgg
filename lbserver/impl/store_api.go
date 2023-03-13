package impl

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"
)

// FileParas 需要请求的content-type为：multipart/form-data
type FileParas struct {
	F *multipart.FileHeader `form:"file"`
}

type UploadRsp struct {
	Url string `json:"url"`
}

func registerStoreApi(h *gin.Engine) {
	h.POST("/upload", Upload)
}

func Upload(c *gin.Context) {
	var req FileParas
	handler := NewHandler(c)

	err := handler.C.ShouldBind(&req)
	if err != nil {
		log.Errorf("err is : %v", err)
		handler.RespError(err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	objectKey := `public/link-info/assets/images/` + req.F.Filename

	open, err := req.F.Open()
	if err != nil {
		log.Errorf("err is %v", err)
		handler.RespError(err)
		return
	}

	err = lb.Storage.Put(objectKey, open)
	if err != nil {
		log.Errorf("err is %v", err)
		handler.RespError(err)
		return
	}

	signedURL, err := lb.Storage.SignURL(objectKey, http.MethodGet, 60*60*24*365)
	if err != nil {
		log.Errorf("err is %v", err)
		handler.RespError(err)
		return
	}

	var rsp = UploadRsp{
		Url: signedURL,
	}
	handler.Success(&rsp)
}

// ConvertMediaUrl 将URL转换为系统的URL
func ConvertMediaUrl(filename, url string) (string, error) {
	objectKey := `public/link-info/assets/images/` + filename

	open, err := http.Get(url)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	err = lb.Storage.Put(objectKey, open.Body)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	signedURL, err := lb.Storage.SignURL(objectKey, http.MethodGet, 60*60*24*365)
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
	err := lb.Storage.Put(objectKey, buffer)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	signedURL, err := lb.Storage.SignURL(objectKey, http.MethodGet, 60*60*24*365)
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}

	return signedURL, nil
}
