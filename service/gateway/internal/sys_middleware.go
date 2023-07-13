package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/lbserver/impl/constant"
	"github.com/oldbai555/bgg/service/gateway/internal/conf"

	"github.com/oldbai555/bgg/pkg/gin_tool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"net/http"
	"strings"
	"time"
)

// Cors 跨域配制
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		method := c.Request.Method
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}

// RegisterUuidTrace 注册一个链路Id进入日志中
func RegisterUuidTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var traceId string
		hint := c.Value(constant.LogWithHint)
		if hint == nil {
			traceId = utils.GenUUID()
		} else if fmt.Sprintf("%v", hint) == "" {
			traceId = utils.GenUUID()
		} else {
			traceId = fmt.Sprintf("%v.%s", hint, utils.GenUUID())
		}
		log.SetLogHint(traceId)
		c.Set(constant.LogWithHint, traceId)
		log.Infof("RemoteIP: %s , ClientIP: %s", c.RemoteIP(), c.ClientIP())
		c.Next()
	}
}

// RegisterJwt 是我们用来检查令牌是否有效的中间件。如果返回401状态无效，则返回给客户。
func RegisterJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cm *gin_tool.Cmd
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

		c.Next()
	}
}

var defaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	hint := param.Keys[constant.LogWithHint]
	return fmt.Sprintf("[GIN] %s<%s> %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		conf.Global.ServerConf.Name,
		hint,
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
