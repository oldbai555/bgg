package impl

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/client/lbim"
	"github.com/oldbai555/lbtool/log"
)

func registerLbimApi(h *gin.Engine) {
	// 可以利用反射来映射函数进去
	group := h.Group("im").Use(RegisterJwt())
}
