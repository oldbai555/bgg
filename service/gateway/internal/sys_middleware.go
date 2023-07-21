package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/pkg/_const"
	_cmd "github.com/oldbai555/bgg/pkg/cmd"
	"github.com/oldbai555/bgg/pkg/gin_tool"
	"net/http"
	"strings"
)

// RegisterJwt 是我们用来检查令牌是否有效的中间件。如果返回401状态无效，则返回给客户。
func RegisterJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cm *_cmd.Cmd
		handler := gin_tool.NewHandler(c)
		for _, cmd := range CmdList {
			if strings.HasSuffix(c.Request.RequestURI, cmd.Path) {
				cm = &cmd
				break
			}
		}

		// 找不到
		if cm == nil {
			// 404
			handler.Response(http.StatusNotFound, http.StatusNotFound, "", "not found")
			return
		}

		// 不用校验权限 - pub
		if !cm.IsUserAuthType() {
			c.Next()
			return
		}

		// todo
		c.Request.Header.Add(_const.HeaderLBSid, "123456789")
		c.Request.Header.Add(_const.HeaderLBCallFrom, cm.GetAuthType())

		c.Next()
	}
}

func RegisterSvr() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("svr", "gateway")
		c.Next()
	}
}
