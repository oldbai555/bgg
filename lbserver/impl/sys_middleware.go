package impl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/bgg/pkg/webtool"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"github.com/storyicon/grbac"
	"net/http"
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
		log.SetModuleName(lb.Sc.Name)
		var traceId string
		hint := c.Value(LogWithHint)
		if hint == nil {
			traceId = utils.GenUUID()
		} else if fmt.Sprintf("%v", hint) == "" {
			traceId = utils.GenUUID()
		} else {
			traceId = fmt.Sprintf("%v.%s", hint, utils.GenUUID())
		}
		log.SetLogHint(traceId)
		c.Set(LogWithHint, traceId)
		c.Next()
	}
}

func RegisterShowReq() gin.HandlerFunc {
	return func(c *gin.Context) {
		mapV2, err := utils.StructToMapV2(c.Request)
		log.Infof("req is %v ,err is %v", mapV2, err)
		c.Next()
	}
}

// RegisterJwt 是我们用来检查令牌是否有效的中间件。如果返回401状态无效，则返回给客户。
func RegisterJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := NewHandler(c)

		// step 1: 拿到 sid
		sid := handler.GetHeader(HttpHeaderAuthorization)

		// step 2: 拿到 token
		token, err := lb.Rdb.Get(c, sid).Result()
		if err != nil {
			log.Errorf("err is : %v", err)
			handler.Response(http.StatusOK, http.StatusUnauthorized, nil, "权限不足")
			c.Abort()
			return
		}

		// step 3: 校验 token
		// vcalidate token formate
		if len(token) == 0 {
			handler.Response(http.StatusOK, http.StatusUnauthorized, nil, "权限不足")
			c.Abort()
			return
		}

		// step 4: 解析 token
		parseToken, claims, err := webtool.ParseToken(token)
		if err != nil || !parseToken.Valid {
			handler.Response(http.StatusOK, http.StatusUnauthorized, nil, "权限不足")
			c.Abort()
			return
		}

		// step 5: 拿到用户所属角色
		roles, err := GetUserRoles(c, claims.UserId)
		if err != nil {
			log.Errorf("err is : %v", err)
			handler.Response(http.StatusOK, http.StatusUnauthorized, nil, "权限不足")
			c.Abort()
			return
		}

		// step 6: 权限校验
		granted, err := lb.Rbac.IsQueryGranted(&grbac.Query{
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Host:   c.Request.Host,
		}, roles)
		if err != nil {
			log.Errorf("err is : %v", err)
			handler.Response(http.StatusOK, http.StatusUnauthorized, nil, "权限不足")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// step 7: 权限校验是否通过
		if !granted.IsGranted() {
			log.Errorf("err is : %s", "权限不足")
			handler.Response(http.StatusOK, http.StatusUnauthorized, nil, "权限不足")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// step 8: 将claims信息放入 context.Context 中
		c.Set(CtxWithClaim, claims)
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

	hint := param.Keys[LogWithHint]
	return fmt.Sprintf("[GIN] %s<%s> %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		lb.Sc.Name,
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
