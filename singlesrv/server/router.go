/**
 * @Author: zjj
 * @Date: 2024/6/17
 * @Desc:
**/

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/singlesrv/client"
	"github.com/oldbai555/bgg/singlesrv/server/wsmgr"
)

// RegisterCustomRouter 注册自定义路由
func RegisterCustomRouter(r *gin.Engine) {
	group := r.Group(client.ServerName)
	group.POST("/upload", handleUploadFile)
	group.GET("/download/*sUrl", handleDownloadFile)
	group.Any("/ws", wsmgr.HandleWs)
}
