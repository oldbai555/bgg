/**
 * @Author: zjj
 * @Date: 2024/5/7
 * @Desc:
**/

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/utils"
	"net/http"
	"path"
)

func registerRouter(router *gin.Engine) {
	router.LoadHTMLGlob(path.Join(utils.GetCurDir(), "templates", "*"))
	router.GET("file.html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "file.html", nil)
	})
	router.POST("/lboss/upload", handleUpload)
	router.POST("/lboss/download", handleDownload)
}

func handleUpload(ctx *gin.Context) {

}

func handleDownload(ctx *gin.Context) {

}
