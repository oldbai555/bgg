/**
 * @Author: zjj
 * @Date: 2024/6/17
 * @Desc:
**/

package lbossserver

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/service/lboss"
)

// RegisterCustomRouter 注册自定义路由
func RegisterCustomRouter(r *gin.Engine) {
	group := r.Group(lboss.ServerName)
	group.POST("/upload", handleUploadFile)
	group.POST("/deploy/upload", handleUploadDeployFile)
	group.GET("/download/*path", handleDownloadFile)
	group.GET("/presigned/*fileName", handlePreSigned)
}
