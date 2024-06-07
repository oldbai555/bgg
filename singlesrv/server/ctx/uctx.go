/**
 * @Author: zjj
 * @Date: 2024/6/7
 * @Desc:
**/

package ctx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldbai555/lbtool/utils"
	"github.com/oldbai555/micro/bcmd"
	"github.com/oldbai555/micro/bconst"
	"github.com/oldbai555/micro/uctx"
)

func NewCtx(ctx context.Context, options ...Option) *Ctx {
	c := &Ctx{
		Context:  ctx,
		BaseUCtx: uctx.NewBaseUCtx(),
	}
	for _, option := range options {
		option(c)
	}
	return c
}

type Ctx struct {
	context.Context
	*uctx.BaseUCtx
}

type Option func(ctx *Ctx)

func WithProtocolType(ctx *gin.Context) Option {
	return func(c *Ctx) {
		val := ctx.GetHeader(bconst.ProtocolType)
		if val != "" {
			c.SetProtocolType(val)
		} else {
			c.SetProtocolType(bconst.PROTO_TYPE_API_JSON) // 默认 json
		}
	}
}

func WithGinHeaderTraceId(ctx *gin.Context) Option {
	return func(c *Ctx) {
		val := ctx.GetHeader(bconst.GinHeaderTraceId)
		if val != "" {
			val = fmt.Sprintf("%s.%s", val, utils.GenRandomStr())
		} else {
			val = utils.GenRandomStr()
		}
	}
}

func WithGinHeaderDeviceId(ctx *gin.Context) Option {
	return func(c *Ctx) {
		val := ctx.GetHeader(bconst.GinHeaderDeviceId)
		if val != "" {
			c.SetDeviceId(val)
		}
	}
}

func WithGinHeaderSid(ctx *gin.Context) Option {
	return func(c *Ctx) {
		val := ctx.GetHeader(bconst.GinHeaderSid)
		if val != "" {
			c.SetSid(val)
		}
	}
}
func WithGinHeaderAuthorization(ctx *gin.Context) Option {
	return func(c *Ctx) {
		val := ctx.GetHeader("Authorization")
		if val != "" {
			c.SetSid(val)
		}
	}
}

func WithGinHeaderAuthType(ctx *gin.Context, cmd *bcmd.Cmd) Option {
	return func(c *Ctx) {
		val := ctx.GetHeader(bconst.GinHeaderAuthType)
		if val != "" {
			c.SetAuthType(val)
		} else {
			c.SetAuthType(cmd.GetAuthType())
		}
	}
}
