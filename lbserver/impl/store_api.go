package impl

import (
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
