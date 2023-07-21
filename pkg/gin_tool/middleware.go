package gin_tool

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/pkg/_const"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"net/http"
	"time"
)

// Cors 跨域配制
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET") // 只放行 POST GET
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
		hint := c.Value(_const.LogWithHint)
		if hint == nil || fmt.Sprintf("%v", hint) == "" {
			hint = c.GetHeader(_const.HeaderLBTraceId)
		}

		if hint == "" {
			traceId = utils.GenUUID()
		} else {
			traceId = fmt.Sprintf("%v.%s", hint, utils.GenUUID())
		}
		log.SetLogHint(traceId)

		c.Set(_const.LogWithHint, traceId)
		c.Request.Header.Add(_const.HeaderLBTraceId, traceId)

		log.Infof("hint: %s , RemoteIP: %s , ClientIP: %s", traceId, c.RemoteIP(), c.ClientIP())

		c.Next()
	}
}

var DefaultLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	hint := param.Keys[_const.LogWithHint]
	svr := param.Keys["svr"]
	v := fmt.Sprintf("[GIN] %s<%s> %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		svr,
		hint,
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)

	_, err := log.GetWriter().Write([]byte(v))
	if err != nil {
		log.Errorf("err:%v", err)
	}

	return v
}
