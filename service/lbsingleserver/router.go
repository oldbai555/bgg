/**
 * @Author: zjj
 * @Date: 2024/6/17
 * @Desc:
**/

package lbsingleserver

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/service/lbsingle"
	"github.com/oldbai555/bgg/service/lbsingleserver/wsmgr"
)

// RegisterCustomRouter 注册自定义路由
func RegisterCustomRouter(r *gin.Engine) {
	group := r.Group(lbsingle.ServerName)
	group.Any("/ws", wsmgr.HandleWs)
}
